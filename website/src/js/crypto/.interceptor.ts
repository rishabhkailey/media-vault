// this package is not being used, it did not work with video and image tag requests
import xhook from "xhook";

export function overrideFetch() {
  // @ts-ignore
  if (window.fetchOverrided === true) {
    return;
  }
  const { fetch: origFetch } = window;
  window.fetch = async (...args) => {
    console.log("fetch called with args:", args);
    const response = await origFetch(...args);

    /* work with the cloned response in a separate promise
       chain -- could use the same chain with `await`. */
    response
      .clone()
      .json()
      .then((body) => console.log("intercepted response:", body))
      .catch((err) => console.error(err));

    /* the original response can be resolved unmodified: */
    //return response;
    return new Response(response.body, {
      headers: response.headers,
      status: response.status,
      statusText: response.statusText,
    });
  };
  // @ts-ignore
  window.fetchOverrided = true;
}

export function overrideXhr() {
  // @ts-ignore
  if (window.xhrOverrided === true) {
    return;
  }
  xhook.after(function (request, response) {
    console.log(response.text);
    // if (request.url.match(/example\.txt$/))
    //   response.text = response.text.replace(/[aeiou]/g, "z");
  });
  /*
    custom solutions - 
    https://stackoverflow.com/questions/16959359/intercept-xmlhttprequest-and-modify-responsetext
    https://stackoverflow.com/questions/45425169/intercept-fetch-api-requests-and-responses-in-javascript
  */
  // @ts-ignore
  window.xhrOverrided = true;
}
