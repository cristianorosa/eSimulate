describe('Fluxo de Login - eSimulate', () => {
  it('deve autenticar com credenciais válidas e redirecionar para o dashboard', () => {
    cy.visit('/login');
    cy.get('input[name=email]').type('user@teste.com');
    cy.get('input[name=password]').type('123456');
    cy.get('button[type=submit]').click();
    cy.contains('Login realizado com sucesso!').should('be.visible');
    cy.url().should('include', '/dashboard');
  });

  it('deve exibir erro com credenciais inválidas', () => {
    cy.visit('/login');
    cy.get('input[name=email]').type('user@teste.com');
    cy.get('input[name=password]').type('errada');
    cy.get('button[type=submit]').click();
    cy.contains('Usuário ou senha inválidos').should('be.visible');
    cy.url().should('include', '/login');
  });
});
