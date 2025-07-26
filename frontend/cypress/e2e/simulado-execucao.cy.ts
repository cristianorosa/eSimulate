describe('Fluxo de Execução de Simulado - eSimulate', () => {
  beforeEach(() => {
    // Autentica antes de cada teste (ajuste conforme necessário para o seu sistema)
    cy.visit('/login');
    cy.get('input[name=email]').type('user@teste.com');
    cy.get('input[name=password]').type('123456');
    cy.get('button[type=submit]').click();
    cy.url().should('include', '/dashboard');
  });

  it('deve executar um simulado e visualizar o resultado', () => {
    cy.visit('/simulados');
    cy.get('button').contains('Iniciar').first().click();
    cy.url().should('include', '/simulados/');
    cy.get('.questao-atual').should('be.visible');
    // Responde todas as questões (seleciona a primeira opção de cada)
    cy.get('mat-radio-group').each(($group) => {
      cy.wrap($group).find('mat-radio-button').first().click();
      cy.get('button').contains('Próxima').click({ force: true });
    });
    cy.get('button').contains('Finalizar').click();
    cy.contains('Resultado do Simulado').should('be.visible');
    cy.contains('Acertos:').should('be.visible');
    cy.contains('Desempenho:').should('be.visible');
  });
});
