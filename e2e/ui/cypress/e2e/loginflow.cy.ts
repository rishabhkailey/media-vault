/// <reference types="cypress" />

// function getVuetifyButtonContaing(
//   contain: string,
// ): Cypress.Chainable<JQuery<HTMLElement>> {
//   return cy
//     .get(".v-btn")
//     .children(".v-btn__content")
//     .contains(contain)
//     .parent();
// }

// new user login
//  * redirect to onboarding
//  * after onboarding no further login attempts should redirect to onboarding
//  * check invalid username and passwords
//  * check invalid access (without media-vault/user role)

describe("login flow", () => {
  beforeEach(() => {
    cy.visit("http://localhost:8181/");
    cy.clearAllCookies();
    cy.clearAllLocalStorage();
    cy.clearAllSessionStorage();
  });

  it("redirct to about page", () => {
    cy.url().should("include", "/about");
  });

  it("sign in > accept consents > correct encryption key", () => {
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
  });
});
