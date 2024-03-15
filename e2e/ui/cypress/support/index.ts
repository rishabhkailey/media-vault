declare namespace Cypress {
  interface Chainable<Subject = any> {
    login(userName: string, password: string): Chainable<any>;
    onBoardUser(encryptionKey: string): Chainable<any>
    enteryEncryptionKey(encryptionKey: string): Chainable<any>
    isMobileVersion(): Chainable<boolean>;
    openSideBarIfMobileVersion(): Chainable<any>;
    closeSideBarIfMobileVersion(): Chainable<any>;
    createAlbum(albumName: string): Chainable<Album>;
    goToHomePage(): Chainable<any>;
    goToAlbumsPage(): Chainable<any>;
    goToAlbumPage(albumId: number): Chainable<any>;
    addMediaToAlbum(albumId: number, mediaIds: Array<number>): Chainable<any>;
    removeMediaFromAlbum(albumId: number, mediaIds: Array<number>): Chainable<any>;
    uploadFiles(files: Array<string>): Chainable<Array<Media>>;
    verifyAlbumsMediaWithSort(albumId: number, mediaList: Array<Media>): Chainable<any>;
    verifyMediaWithSort(mediaList: Array<Media>): Chainable<Array<any>>;
    selectMedia(mediaIds: Array<number>): Chainable<any>;
    deleteAlbum(albumId: number): Chainable<any>;
    deleteMedia(mediaIds: Array<number>): Chainable<any>;
    verifyMediaCarausel(mediaList: Array<Media>, indexOffSet: number, urlPathPrefix: string): Chainable<any>;
  }
}