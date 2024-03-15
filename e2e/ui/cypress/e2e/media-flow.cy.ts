/// <reference types="cypress" />
import { sortMediaByUploadDateDesc } from "./utils";

// for intializing the post login test scenarios
function loginWithTestUser() {
  // todo skip if already logged in
  cy.clearAllCookies();
  cy.clearAllLocalStorage();
  cy.clearAllSessionStorage();

  cy.visit("http://localhost:8181/");

  cy.get(`[data-test-id="signin-button"]`).click();
  // todo realm name can change
  cy.url().should("include", "/accounts/realms").should("include", "/protocol/openid-connect/auth");

  cy.get("#username").type("test-user");
  cy.get("#password").type("ZmuBYNurSS");
  cy.get("#kc-login").click();

  cy.url().should("contain", "/encryption-key");
  cy.get(`[data-test-id="encryption-key-input"]`).type("ZmuBYNurSS");
  cy.get(`[data-test-id="encryption-key-submit-button"]`).click();
  cy.url().should("not.contain", "/encryption-key");
  cy.url().should("not.contain", "/about");
  cy.url().should("not.contain", "/realms");
}

// const testFiles = [
//   "test-files/image-01.png",
//   "test-files/image-02.png",
//   "test-files/image-03.png",
//   "test-files/video-01.webm",
// ];

describe("media flow", () => {
  beforeEach(() => {
    cy.login(Cypress.env("TEST_USER_NAME"), Cypress.env("TEST_USER_PASSWORD"))
    cy.enteryEncryptionKey(Cypress.env("TEST_USER_ENCRYPTION_KEY"))
    loginWithTestUser();
    // cy.viewport(1280, 1080)
  });
  it("upload files", () => {
    let uploadedMedia: Array<Media> = []
    let createdAlbum: Album;
    cy.fixture<Array<string>>("test-files").then((testFiles) => {
      cy.uploadFiles(testFiles).then((_uploadedMedia) => {
        uploadedMedia = _uploadedMedia
        cy.verifyMediaWithSort(sortMediaByUploadDateDesc(uploadedMedia))
      })
      cy.createAlbum("test-album-01").then((album) => {
        createdAlbum = album
      }).then(() => {
        cy.log(JSON.stringify(uploadedMedia))
        cy.addMediaToAlbum(createdAlbum.id, sortMediaByUploadDateDesc(uploadedMedia).map(m => m.id))
        cy.verifyAlbumsMediaWithSort(createdAlbum.id, uploadedMedia)
      }).then(() => {
        cy.removeMediaFromAlbum(createdAlbum.id, uploadedMedia.map(m => m.id))
      }).then(() => cy.deleteAlbum(createdAlbum.id)).then(() => {
        cy.deleteMedia(uploadedMedia.map(m => m.id))
      })
    })

    // search by exact file name
    // search by type image/video
    // search by date
  });
})