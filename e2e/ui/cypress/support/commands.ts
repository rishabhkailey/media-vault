/// <reference types="cypress" />

import cypress = require("cypress");

function getBaseFileName(filePath: string): string {
  if (
    filePath.indexOf("/") !== -1 ||
    filePath.indexOf("/") !== filePath.length - 1
  ) {
    return filePath.substring(filePath.indexOf("/") + 1);
  }
  return filePath;
}

Cypress.Commands.add("login", (userName: string, password: string) => {
  cy.clearAllCookies();
  cy.clearAllLocalStorage();
  cy.clearAllSessionStorage();

  cy.visit("/");

  cy.get(`[data-test-id="signin-button"]`).click();
  cy.url().should("include", "/accounts/realms").should("include", "/protocol/openid-connect/auth");

  cy.get("#username").type(userName);
  cy.get("#password").type(password);
  cy.get("#kc-login").click();
  cy.url().should("not.include", "/accounts")
  return cy.url().should("not.include", "/error")
})

Cypress.Commands.add("onBoardUser", (encryptionKey: string) => {
  // todo
  return cy.wait(0)
})

Cypress.Commands.add("enteryEncryptionKey", (encryptionKey: string) => {
  cy.url().should("contain", "/encryption-key");
  cy.get(`[data-test-id="encryption-key-input"]`).type(encryptionKey);
  cy.get(`[data-test-id="encryption-key-submit-button"]`).click();
  cy.url().should("not.contain", "/encryption-key");
  cy.url().should("not.contain", "/about");
  return cy.url().should("not.contain", "/realms");
})

Cypress.Commands.add("isMobileVersion", () => {
  return cy.document().then(($document) => cy.wait(100).then(() => {
    const documentResult = $document.querySelectorAll(`[data-test-id="open-close-sidebar-button"]`)
    if (documentResult.length > 0) {
      return true
    }
    return false
  }))
})

Cypress.Commands.add("openSideBarIfMobileVersion", () => {
  return cy.isMobileVersion().then((isMobile) => {
    if (!isMobile) {
      return false
    }
    return cy.get(`[data-test-id=side-bar]`).then(($sideBar) => {
      cy.log($sideBar.hasClass('v-navigation-drawer--active') + ' - visible side bar')
      if ($sideBar.hasClass('v-navigation-drawer--active')) {
        return
      }
      cy.get(`[data-test-id="open-close-sidebar-button"]`).click()
      return
    })
  })
})

Cypress.Commands.add("closeSideBarIfMobileVersion", () => {
  return cy.isMobileVersion().then((isMobile) => {
    if (!isMobile) {
      return false
    }
    return cy.get(`[data-test-id=side-bar]`).then(($sideBar) => {
      cy.log($sideBar.hasClass('v-navigation-drawer--active') + ' - visible side bar')
      if ($sideBar.hasClass('v-navigation-drawer--active')) {
        cy.get(`[data-test-id="open-close-sidebar-button"]`).click()
      }
      return
    })
  })
})

Cypress.Commands.add("createAlbum", (albumName: string) => {
  cy.goToAlbumsPage()
  cy.intercept(
    {
      method: "POST",
      path: "/v1/album",
    },
  ).as("createAlbumRequest")
  cy.get(`[data-test-id="album-create-button"]`).click()
  cy.get(`[data-test-id="create-album-form-name-input"] input`).type(albumName)
  cy.get(`[data-test-id="create-album-form-submit-button"]`).click()
  return cy.wait("@createAlbumRequest", {
    timeout: 10000,
  }).then((interception) => {
    if (interception.response.statusCode !== 200) {
      throw new Error(`create album failed with status code ${interception.response.statusCode}`)
    }
    const albumResponse = interception.response.body as Album
    cy.log(JSON.stringify(albumResponse))
    return cy.get(`[data-test-id="album_card_container_${albumResponse.id}"]`).should("exist").then(() => {
      return albumResponse
    })
  });
})

Cypress.Commands.add("goToHomePage", () => {
  return cy.openSideBarIfMobileVersion().then(() => {
    return cy.get(`[data-test-id="sidebar-home-group"]`).click().then(() => cy.closeSideBarIfMobileVersion())
  })
})

Cypress.Commands.add("goToAlbumsPage", () => {
  cy.openSideBarIfMobileVersion()
  cy.get(`[data-test-id="sidebar-albums-group"]`).click().then(() => cy.closeSideBarIfMobileVersion())
  return cy.url().should("include", "/albums")
})

Cypress.Commands.add("selectMedia", (mediaIds: Array<number>) => {
  mediaIds.forEach((mediaId) => {
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).should("exist")
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).trigger("mouseover")
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).trigger("mouseenter")
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"] [data-test-id="select-button"]`).should("be.visible")
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"] [data-test-id="select-button"]`).click()
  })
  return cy.wait(0)
})

Cypress.Commands.add("addMediaToAlbum", (albumId: number, mediaIds: Array<number>) => {
  cy.intercept(
    {
      method: "POST",
      path: "/v1/album/**/media",
    },
  ).as("addMediaToAlbumRequest")

  cy.goToHomePage()
  cy.selectMedia(mediaIds)
  cy.get(`[data-test-id="appbar-add-to-album-button"]`).click()
  cy.get(`[data-test-id="appbar-add-to-album-checkbox-${albumId}"] input[type="checkbox"]`).click()
  cy.get(`[data-test-id="appbar-add-to-album-confirm-button"]`).click()
  cy.goToAlbumPage(albumId)

  mediaIds.forEach((mediaId) => {
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).should("exist")
  })

  return cy.wait("@addMediaToAlbumRequest", {
    timeout: 10000,
  }).then((interception) => {
    if (interception.response.statusCode !== 200) {
      throw new Error(`add media to album request failed with status code ${interception.response.statusCode}`)
    }
    return true
  })
})

Cypress.Commands.add("removeMediaFromAlbum", (albumId: number, mediaIds: Array<number>) => {
  cy.intercept(
    {
      method: "DELETE",
      path: `/v1/album/${albumId}/media`,
    },
  ).as("removeMediaFromAlbumRequest")

  cy.goToAlbumPage(albumId)
  cy.selectMedia(mediaIds)
  cy.get(`button[data-test-id="appbar-album-actions-remove-from-album"]`).click()
  cy.get(`[data-test-id="appbar-album-actions-remove-from-album-confrim"] button[data-test-id="confirm-button"]`).click()

  mediaIds.forEach((mediaId) => {
    cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).should("not.exist")
  })

  return cy.wait("@removeMediaFromAlbumRequest", {
    timeout: 10000,
  }).then((interception) => {
    if (interception.response.statusCode !== 200) {
      throw new Error(`remove media from album request failed with status code ${interception.response.statusCode}`)
    }
    return true
  }).then(() => {
    mediaIds.forEach((mediaId) => {
      cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).should("not.exist")
    })
  })
})

Cypress.Commands.add("uploadFiles", (files: Array<string>) => {
  cy.intercept({
    method: "POST",
    path: "/v1/upload",
  }).as("initUploadRequests");

  cy.intercept(
    {
      method: "POST",
      path: "/v1/upload/**/finish",
    },
  ).as("finishUploadRequest");

  // select files and upload
  cy.get(`[data-test-id="upload-file-button"]`).click();
  cy.get(`[data-test-id="upload-file-input"] input[type="file"]`).selectFile(
    files,
  );
  cy.get(`[data-test-id="upload-files-selection-confirm-button"]`).click();

  // wait for finishUploadRequests
  const uploadedMedia: Array<Media> = [];
  cy.wait(Array(files.length).fill("@finishUploadRequest"), {
    timeout: 10000,
  })
    .then((interceptions) => {
      cy.log(
        interceptions.length.toString(),
        interceptions[0].response.statusMessage,
      );
      interceptions.forEach((interception) => {
        cy.log(interception.response.statusCode.toString());
        if (interception.response.statusCode === 200) {
          uploadedMedia.push(interception.response.body as Media)
        }
      });
    }).then(() => {
      expect(uploadedMedia.length).equal(files.length)
    })

  // check upload files dialog
  cy.get(`[data-test-id="uploading-files-dialog-notcallapsed"]`).should("exist");
  cy.get(`[data-test-id="uploading-files-progress-list"]`).should("exist");
  cy.get(
    `[data-test-id="uploading-files-progress-list"] .v-list-item__content .v-list-item-title`,
  )
    .should("have.length", files.length)
    .each(($element, index) => {
      const fileName = getBaseFileName(files[index]);
      cy.log(
        `element text = ${$element.text()} expected file name = ${fileName}`,
      );
      cy.wrap($element).contains(fileName);
    });

  return cy.get(
    `[data-test-id="uploading-files-progress-list"] .v-list-item__prepend > .v-avatar > .v-progress-circular > .v-progress-circular__content > span`,
  )
    .should("have.length", files.length)
    .each(($element) => {
      cy.wrap($element).contains("100%");
    }).then(() => uploadedMedia)
})

Cypress.Commands.add("goToAlbumPage", (albumId: number) => {
  cy.goToAlbumsPage()
  cy.get(`[data-test-id="album_card_container_${albumId}"]`).should("exist")
  return cy.get(`[data-test-id="album_card_container_${albumId}"] .v-card--link`).click()
})

// check media from start
Cypress.Commands.add("verifyAlbumsMediaWithSort", (albumId: number, mediaList: Array<Media>) => {
  cy.goToAlbumPage(albumId)
  const thumbnailMediaIds: Array<number> = []
  cy.get(`[data-test-id*="thumbnail_container_"]`).each(($thumbnailContainer) => {
    thumbnailMediaIds.push(Number($thumbnailContainer.attr('data-test-id').replace('thumbnail_container_', '')))
  }).then(() => {
    cy.log(JSON.stringify(thumbnailMediaIds), JSON.stringify(mediaList.map<number>(m => m.id)))
    expect(thumbnailMediaIds).to.deep.eq(mediaList.map<number>(m => m.id))
  })
  cy.verifyMediaCarausel(mediaList, 0, `/album/${albumId}`)
})

Cypress.Commands.add("verifyMediaCarausel", (mediaList: Array<Media>, indexOffSet: number, urlPathPrefix: string) => {
  mediaList.forEach((media, index) => {
    // caching of images/video by browser causes issue if we test media carausel more than 1 time (e.g. normal media carausel, album media carausel)
    // once the media is opened by browser once it doesn't send a new request on next time we open the media
    // const mediaRequestName = `media_file_request_${media.id}_${media.url}`.replace("/", "_")
    // cy.intercept(
    //   {
    //     method: "GET",
    //     path: media.url,
    //   }, (request) => {
    //     // request.headers['Cache-Control'] = 'no-cache';
    //     request.continue((responseInterceptor) => {
    //       // to prevent caching of images/video as it causes issue if we test media carausel more than 1 time (e.g. normal media carausel, album media carausel)
    //       responseInterceptor.headers['Cache-Control'] = 'no-cache';
    //     });
    //   }
    // ).as(mediaRequestName)
    
    if (index === 0) {
      // open media carausel if first media
      cy.get(`[data-test-id*="thumbnail_container_${media.id}"]`).first().click()
    } else if (index !== mediaList.length) /* not length - 1, because we are at the previous media, in a loop we either open a media or click next then wait for the request */ {
      // go to next media if not the first media
      cy.get(`.v-window__controls button[aria-label="Next visual"]`).click()
    }
    cy.url().should("include", `${urlPathPrefix}/media/${media.id}/index/${index + indexOffSet}`)
  //   cy.wait(`@${mediaRequestName}`, {
  //     timeout: 10000,
  //   }).then((interception) => {
  //     if (![200, 206].includes(interception.response.statusCode)) {
  //       throw new Error(`get Media File Request failed with status code ${interception.response.statusCode}`)
  //     }
  //   })
  })
  return cy.get(`[data-test-id="media-carousel-close-button"]`).click()
})

Cypress.Commands.add("verifyMediaWithSort", (mediaList: Array<Media>) => {
  cy.log(JSON.stringify(mediaList.map<number>(m => m.id)))
  cy.goToHomePage()
  const thumbnailMediaIds: Array<number> = []
  cy.get(`[data-test-id*="thumbnail_container_"]`).each(($thumbnailContainer) => {
    thumbnailMediaIds.push(Number($thumbnailContainer.attr('data-test-id').replace('thumbnail_container_', '')))
  }).then(() => {
    cy.log(JSON.stringify(thumbnailMediaIds.slice(0, mediaList.length)), JSON.stringify(mediaList.map<number>(m => m.id)))
    expect(thumbnailMediaIds.slice(0, mediaList.length)).to.deep.eq(mediaList.map(m => m.id))
  })

  cy.verifyMediaCarausel(mediaList, 0, "")
})

Cypress.Commands.add("deleteAlbum", (albumId: number) => {
  cy.intercept(
    {
      method: "DELETE",
      path: `/v1/album/${albumId}`,
    },
  ).as("deleteAlbumRequest")

  // todo wait for request to complete
  cy.goToAlbumPage(albumId)
  cy.get(`[data-test-id="delete-album-button"]`).click()
  cy.get(`[data-test-id="delete-album-confirmation"] button[data-test-id="confirm-button"]`).click()

  return cy.wait("@deleteAlbumRequest", {
    timeout: 10000,
  }).then((interception) => {
    if (interception.response.statusCode !== 200) {
      throw new Error(`delete alubm request failed with status code ${interception.response.statusCode}`)
    }
    return true
  })
})

Cypress.Commands.add("deleteMedia", (mediaIds: Array<number>) => {
  cy.intercept(
    {
      method: "DELETE",
      path: `/v1/media`,
    },
  ).as("deleteMediaRequest")

  // todo wait for request to complete
  cy.goToHomePage()
  cy.selectMedia(mediaIds)
  cy.get(`[data-test-id="appbar-delete-media-button"]`).click()
  cy.get(`[data-test-id="appbar-delete-media-confirmation"] button[data-test-id="confirm-button"]`).click()

  return cy.wait("@deleteMediaRequest", {
    timeout: 10000,
  }).then((interception) => {
    if (interception.response.statusCode !== 200) {
      throw new Error(`delete media request faild with status code ${interception.response.statusCode}`)
    }
    return true
  }).then(() => {
    mediaIds.forEach((mediaId) => {
      cy.get(`[data-test-id="thumbnail_container_${mediaId}"]`).should("not.exist")
    })
  })
})