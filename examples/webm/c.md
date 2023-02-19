[source](https://github.com/FFmpeg/FFmpeg/blob/master/doc/examples/transcoding.c)

```c
// open_input_file
// seems to be just intializing codec ctx from the opened file
static AVFormatContext *ifmt_ctx;
avformat_open_input(&ifmt_ctx, filename, NULL, NULL); // < 0 means errro
avformat_find_stream_info(ifmt_ctx, NULL) // update stream info in ifmt_ctx by reading the file
// after this ifmt-ctx->nb_streams = no. of streams
// ifmt_ctx->streams[i] = array of streams

// this is global
typedef struct StreamContext {
    // input stream codec ctx
    AVCodecContext *dec_ctx;
    AVCodecContext *enc_ctx;

    AVFrame *dec_frame;
} StreamContext;
static StreamContext *stream_ctx;
// array lenght = no of streams in the input file
stream_ctx = av_calloc(ifmt_ctx->nb_streams, sizeof(*stream_ctx));
for (i = 0; i < ifmt_ctx->nb_streams; i++) {
  AVStream *stream = ifmt_ctx->streams[i];
  // get stream codec
  const AVCodec *dec = avcodec_find_decoder(stream->codecpar->codec_id);
  // codec ctx
  codec_ctx = avcodec_alloc_context3(dec);
  // copy stream params to codec ctx (format, height, width, sample rate) - check here https://ffmpeg.org/doxygen/trunk/structAVCodecParameters.html
  avcodec_parameters_to_context(codec_ctx, stream->codecpar);
  if (codec_ctx->codec_type == AVMEDIA_TYPE_VIDEO || codec_ctx->codec_type == AVMEDIA_TYPE_AUDIO) {
    if (codec_ctx->codec_type == AVMEDIA_TYPE_VIDEO)
        // if video then set framrate of codec ctx
        codec_ctx->framerate = av_guess_frame_rate(ifmt_ctx, stream, NULL);
    /* Open decoder */
    ret = avcodec_open2(codec_ctx, dec, NULL);
    if (ret < 0) {
        av_log(NULL, AV_LOG_ERROR, "Failed to open decoder for stream #%u\n", i);
        return ret;
    }
    // set the input codec ctx
    stream_ctx[i].dec_ctx = codec_ctx;
    stream_ctx[i].dec_frame = av_frame_alloc();
  }
}
```

```c
// open output file
static AVFormatContext *ofmt_ctx;
avformat_alloc_output_context2(&ofmt_ctx, NULL, NULL, filename);
// for every input stream we will create 1 output stream
for (i = 0; i < ifmt_ctx->nb_streams; i++) {
  // create output stream 
  out_stream = avformat_new_stream(ofmt_ctx, NULL);

  in_stream = ifmt_ctx->streams[i];
  // decode/input codec ctx we intialized in open input file
  dec_ctx = stream_ctx[i].dec_ctx;
  if (dec_ctx->codec_type == AVMEDIA_TYPE_VIDEO
        || dec_ctx->codec_type == AVMEDIA_TYPE_AUDIO) {
    /* in this example, we choose transcoding to same codec */
    encoder = avcodec_find_encoder(dec_ctx->codec_id);
    enc_ctx = avcodec_alloc_context3(encoder);
    /* In this example, we transcode to same properties (picture size,
      * sample rate etc.). These properties can be changed for output
      * streams easily using filters */
    if (dec_ctx->codec_type == AVMEDIA_TYPE_VIDEO) {
        enc_ctx->height = dec_ctx->height;
        enc_ctx->width = dec_ctx->width;
        enc_ctx->sample_aspect_ratio = dec_ctx->sample_aspect_ratio;
        /* take first format from list of supported formats */
        if (encoder->pix_fmts)
            enc_ctx->pix_fmt = encoder->pix_fmts[0];
        else
            enc_ctx->pix_fmt = dec_ctx->pix_fmt;
        /* video time_base can be set to whatever is handy and supported by encoder */
        enc_ctx->time_base = av_inv_q(dec_ctx->framerate);
    } else {
        // audio stream
        enc_ctx->sample_rate = dec_ctx->sample_rate;
        ret = av_channel_layout_copy(&enc_ctx->ch_layout, &dec_ctx->ch_layout);
        if (ret < 0)
            return ret;
        /* take first format from list of supported formats */
        enc_ctx->sample_fmt = encoder->sample_fmts[0];
        enc_ctx->time_base = (AVRational){1, enc_ctx->sample_rate};
    }
    if (ofmt_ctx->oformat->flags & AVFMT_GLOBALHEADER)
      enc_ctx->flags |= AV_CODEC_FLAG_GLOBAL_HEADER;
    // open encoder
    ret = avcodec_open2(enc_ctx, encoder, NULL);
    // copy the codec context params to output stream
    ret = avcodec_parameters_from_context(out_stream->codecpar, enc_ctx);
    out_stream->time_base = enc_ctx->time_base;
    stream_ctx[i].enc_ctx = enc_ctx;
  } else if (dec_ctx->codec_type == AVMEDIA_TYPE_UNKNOWN) {
    // error
  } else {
    /* if this stream must be remuxed */
    // remux is for copying the data from 1 container to another
    ret = avcodec_parameters_copy(out_stream->codecpar, in_stream->codecpar);
    if (ret < 0) {
        av_log(NULL, AV_LOG_ERROR, "Copying parameters for stream #%u failed\n", i);
        return ret;
    }
    out_stream->time_base = in_stream->time_base;
  }
}

// log the the details
av_dump_format(ofmt_ctx, 0, filename, 1);

if (!(ofmt_ctx->oformat->flags & AVFMT_NOFILE)) {
    ret = avio_open(&ofmt_ctx->pb, filename, AVIO_FLAG_WRITE);
}
/* init muxer, write output file header */
// Allocate the stream private data and write the stream header to an output media file.
ret = avformat_write_header(ofmt_ctx, NULL);
```

```c
// init_filters

typedef struct FilteringContext {
    AVFilterContext *buffersink_ctx;
    AVFilterContext *buffersrc_ctx;
    AVFilterGraph *filter_graph;

    AVPacket *enc_pkt;
    AVFrame *filtered_frame;
} FilteringContext;
static FilteringContext *filter_ctx;
// array size = no of streams
filter_ctx = av_malloc_array(ifmt_ctx->nb_streams, sizeof(*filter_ctx));
for (i = 0; i < ifmt_ctx->nb_streams; i++) {
  filter_ctx[i].buffersrc_ctx  = NULL;
  filter_ctx[i].buffersink_ctx = NULL;
  filter_ctx[i].filter_graph   = NULL;
  // do nothing if stream is not a audio or video
  if (!(ifmt_ctx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_AUDIO
        || ifmt_ctx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO))
    continue;
  
  // will be used by init filter?
  if (ifmt_ctx->streams[i]->codecpar->codec_type == AVMEDIA_TYPE_VIDEO)
    filter_spec = "null"; /* passthrough (dummy) filter for video */
  else
    filter_spec = "anull"; /* passthrough (dummy) filter for audio */
  
  ret = init_filter(&filter_ctx[i], stream_ctx[i].dec_ctx,
          stream_ctx[i].enc_ctx, filter_spec);
  
```

```c
// init_filter
char args[512];
int ret = 0;
const AVFilter *buffersrc = NULL;
const AVFilter *buffersink = NULL;
AVFilterContext *buffersrc_ctx = NULL;
AVFilterContext *buffersink_ctx = NULL;
AVFilterInOut *outputs = avfilter_inout_alloc();
AVFilterInOut *inputs  = avfilter_inout_alloc();
// create filter graph
AVFilterGraph *filter_graph = avfilter_graph_alloc();

if (dec_ctx->codec_type == AVMEDIA_TYPE_VIDEO) {
  // memory buffer for source api https://ffmpeg.org/doxygen/3.3/group__lavfi__buffersrc.html
  buffersrc = avfilter_get_by_name("buffer");
  // memory buffer for sink api https://ffmpeg.org/doxygen/3.3/group__lavfi__buffersink.html
  buffersink = avfilter_get_by_name("buffersink");
  // similar to fmt.Fprintf
  snprintf(args, sizeof(args),
          "video_size=%dx%d:pix_fmt=%d:time_base=%d/%d:pixel_aspect=%d/%d",
          dec_ctx->width, dec_ctx->height, dec_ctx->pix_fmt,
          dec_ctx->time_base.num, dec_ctx->time_base.den,
          dec_ctx->sample_aspect_ratio.num,
          dec_ctx->sample_aspect_ratio.den);

  // create filter instance in filter_graph, "in" is name
  ret = avfilter_graph_create_filter(&buffersrc_ctx, buffersrc, "in", args, NULL, filter_graph);
  ret = avfilter_graph_create_filter(&buffersink_ctx, buffersink, "out", NULL, NULL, filter_graph);
  // args/opts for buffersink (buffersrc output were set during avfilter_graph_create_filter)
  ret = av_opt_set_bin(buffersink_ctx, "pix_fmts",
          (uint8_t*)&enc_ctx->pix_fmt, sizeof(enc_ctx->pix_fmt),
          AV_OPT_SEARCH_CHILDREN);
} else if (dec_ctx->codec_type == AVMEDIA_TYPE_AUDIO) {
  buffersrc = avfilter_get_by_name("abuffer");
  buffersink = avfilter_get_by_name("abuffersink");

  // if Only the channel count is specified, without any further information about the channel order.
  if (dec_ctx->ch_layout.order == AV_CHANNEL_ORDER_UNSPEC)
    av_channel_layout_default(&dec_ctx->ch_layout, dec_ctx->ch_layout.nb_channels);
  

  char buf[64];
  // Get a human-readable string describing the channel layout properties. in buf
  av_channel_layout_describe(&dec_ctx->ch_layout, buf, sizeof(buf));

  // set args
  snprintf(args, sizeof(args),
        "time_base=%d/%d:sample_rate=%d:sample_fmt=%s:channel_layout=%s",
        dec_ctx->time_base.num, dec_ctx->time_base.den, dec_ctx->sample_rate,
        av_get_sample_fmt_name(dec_ctx->sample_fmt),
        buf);

  // create filter instance in filter_graph, "in" is name
  ret = avfilter_graph_create_filter(&buffersrc_ctx, buffersrc, "in", args, NULL, filter_graph);
  ret = avfilter_graph_create_filter(&buffersink_ctx, buffersink, "out", NULL, NULL, filter_graph);
  
  // args/opts for buffersink (buffersrc output were set during avfilter_graph_create_filter)
  ret = av_opt_set_bin(buffersink_ctx, "sample_fmts",
        (uint8_t*)&enc_ctx->sample_fmt, sizeof(enc_ctx->sample_fmt),
        AV_OPT_SEARCH_CHILDREN);

  av_channel_layout_describe(&enc_ctx->ch_layout, buf, sizeof(buf));
  ret = av_opt_set(buffersink_ctx, "ch_layouts",
          buf, AV_OPT_SEARCH_CHILDREN);
  ret = av_opt_set(buffersink_ctx, "ch_layouts",
          buf, AV_OPT_SEARCH_CHILDREN);
  ret = av_opt_set_bin(buffersink_ctx, "sample_rates",
          (uint8_t*)&enc_ctx->sample_rate, sizeof(enc_ctx->sample_rate),
          AV_OPT_SEARCH_CHILDREN);

}
```


ffmpeg filter examples
```bash
ffmpeg -i INPUT -vf "split [main][tmp]; [tmp] crop=iw:ih/2:0:0, vflip [flip]; [main][flip] overlay=0:H/2" OUTPUT
# split stream into 2 main, tmp
# crop tmp stream iw:ih/2(width, height) :0:0(offset) first upper half and vertically flip the video
# [main][flip]  combine these 2 streams with (overlay flip on main) overlay=0:/H2 (position of overlay)

# change the overlay position to the top (just flip upperhalf not changing lowerhalf)
ffmpeg -i input/test.mp4 -vf "split [main][tmp]; [tmp] crop=iw:ih/2:0:0, vflip [flip]; [main][flip] overlay=0:0" output/filtered.mp4


# in overlay w = width of overlay, W = main video width
# verticly flip middle part of the video
ffmpeg -i input/test.mp4 -vf "split [main][tmp]; [tmp] crop=iw/2:ih/2:iw/4:ih/4, vflip [flip]; [main][flip] overlay=W/4:H/4" output/filtered.mp4

# check more examples here

```