describe('Fluxo de Cadastro de Questão - eSimulate', () => {
  beforeEach(() => {
    // Autentica antes de cada teste (ajuste conforme necessário para o seu sistema)
    cy.visit('/login');
    cy.get('input[name=email]').type('user@teste.com');
    cy.get('input[name=password]').type('123456');
    cy.get('button[type=submit]').click();
    cy.url().should('include', '/dashboard');
  });

  it('deve cadastrar uma nova questão com sucesso', () => {
    cy.visit('/questoes/new');
    cy.get('textarea[name=statement]').type('Enunciado E2E ' + Date.now());
    cy.get('input[name=theme_id]').type('1');
    cy.get('input[name=option0]').type('Opção A');
    cy.get('input[name=option1]').type('Opção B');
    cy.get('input[name=correct0]').check();
    cy.get('button[type=submit]').click();
    cy.contains('Questão cadastrada com sucesso!').should('be.visible');
    cy.url().should('include', '/questoes');
  });
});
