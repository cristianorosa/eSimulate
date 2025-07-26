describe('Fluxo de Cadastro - eSimulate', () => {
  it('deve cadastrar usuário com dados válidos e redirecionar para login', () => {
    cy.visit('/register');
    cy.get('input[name=name]').type('Usuário E2E');
    cy.get('input[name=email]').type('e2e_' + Date.now() + '@teste.com');
    cy.get('input[name=password]').type('123456');
    cy.get('button[type=submit]').click();
    cy.contains('Cadastro realizado com sucesso!').should('be.visible');
    cy.url().should('include', '/login');
  });

  it('deve exibir erro ao tentar cadastrar com email já existente', () => {
    cy.visit('/register');
    cy.get('input[name=name]').type('Usuário Existente');
    cy.get('input[name=email]').type('user@teste.com');
    cy.get('input[name=password]').type('123456');
    cy.get('button[type=submit]').click();
    cy.contains('Erro ao cadastrar').should('be.visible');
    cy.url().should('include', '/register');
  });
});
