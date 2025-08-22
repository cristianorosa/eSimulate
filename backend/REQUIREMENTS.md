# Requisitos do Sistema eSimulate

Documentação completa dos requisitos funcionais e não funcionais do sistema.

## 📋 Visão Geral

O eSimulate é um sistema de simulados educacionais que permite:
- Gestão hierárquica de áreas, exames, tópicos e questões
- Relacionamentos N:N flexíveis entre entidades
- Sistema de tags para categorização de questões
- Importação de questões via JSON
- Validação inteligente de regras de negócio

## 🎯 Requisitos Funcionais

### RF01 - Gestão de Áreas
- **Descrição**: Sistema deve permitir CRUD completo de áreas de conhecimento
- **Critérios de Aceitação**:
  - ✅ Criar área com nome único e descrição opcional
  - ✅ Listar áreas com paginação e busca
  - ✅ Editar informações da área
  - ✅ Excluir área (cascade para exames)
  - ✅ Validar nome único por área

### RF02 - Gestão de Exames
- **Descrição**: Sistema deve permitir criação e gestão de exames
- **Critérios de Aceitação**:
  - ✅ Exame pertence a uma área específica
  - ✅ Configurar tempo máximo (em minutos)
  - ✅ Definir nota mínima para aprovação
  - ✅ Status ativo/inativo
  - ✅ Filtros por área e status
  - ✅ Paginação e busca

### RF03 - Gestão de Tópicos
- **Descrição**: Sistema deve permitir criação de tópicos independentes
- **Critérios de Aceitação**:
  - ✅ Tópico é entidade independente (não pertence a exame específico)
  - ✅ Nome único e descrição opcional
  - ✅ Pode ser associado a múltiplos exames
  - ✅ CRUD completo com validações

### RF04 - Relacionamento Exame-Tópico (N:N)
- **Descrição**: Exames podem ter múltiplos tópicos com configurações específicas
- **Critérios de Aceitação**:
  - ✅ Associação N:N entre exames e tópicos
  - ✅ Cada associação define:
    - Quantidade de questões do tópico no exame
    - Peso percentual na nota final (0-100%)
    - Ordem de apresentação
    - Distribuição de dificuldade (fácil/médio/difícil %)
  - ✅ Soma dos pesos deve ser 100% por exame
  - ✅ Distribuição de dificuldade deve somar 100%

### RF05 - Gestão de Questões
- **Descrição**: Sistema deve permitir criação e gestão de questões
- **Critérios de Aceitação**:
  - ✅ Questão pertence a um tópico específico (1:N)
  - ✅ Texto da questão obrigatório
  - ✅ Nível de dificuldade: fácil, médio, difícil
  - ✅ Explicação opcional
  - ✅ Status ativo/inativo
  - ✅ Múltiplas opções de resposta (mínimo 2)
  - ✅ Apenas uma opção correta por questão
  - ✅ Filtros por exame, tópico, dificuldade

### RF06 - Relacionamento Exame-Questão (N:N)
- **Descrição**: Questões podem pertencer a múltiplos exames com validação hierárquica
- **Critérios de Aceitação**:
  - ✅ Associação N:N entre exames e questões
  - ✅ **REGRA PRINCIPAL**: Questão só pode ser associada a exame se seu tópico estiver associado ao exame
  - ✅ **EXCEÇÃO**: Exame com único tópico (peso 100%) aceita qualquer questão
  - ✅ Ordem de apresentação configurável
  - ✅ Validação implementada na camada de negócio (não no banco)

### RF07 - Sistema de Tags
- **Descrição**: Sistema de categorização flexível para questões
- **Critérios de Aceitação**:
  - ✅ Tag possui apenas nome único (sem descrição)
  - ✅ Relacionamento N:N entre questões e tags
  - ✅ CRUD completo de tags
  - ✅ Busca por nome de tag
  - ✅ Filtro de questões por tags
  - ✅ Associação/desassociação em lote

### RF08 - Importação de Questões
- **Descrição**: Sistema deve permitir importação via JSON
- **Critérios de Aceitação**:
  - ✅ Importar estrutura completa: área → exame → tópico → questões
  - ✅ Criar entidades se não existirem
  - ✅ Atualizar entidades existentes
  - ✅ Associar tags às questões
  - ✅ Validar integridade dos dados
  - ✅ Rollback em caso de erro

### RF09 - Paginação e Filtros
- **Descrição**: Todas as listagens devem suportar paginação e filtros
- **Critérios de Aceitação**:
  - ✅ Paginação padrão: página 1, 10 itens
  - ✅ Filtros específicos por entidade
  - ✅ Busca textual quando aplicável
  - ✅ Contadores de total de itens

### RF10 - Validações de Negócio
- **Descrição**: Sistema deve validar regras de negócio na camada de aplicação
- **Critérios de Aceitação**:
  - ✅ Nomes únicos onde aplicável
  - ✅ Pesos percentuais somam 100%
  - ✅ Distribuição de dificuldade soma 100%
  - ✅ Validação hierárquica exame-questão
  - ✅ Mensagens de erro específicas e úteis

## 🔧 Requisitos Não Funcionais

### RNF01 - Performance
- **Descrição**: Sistema deve ter performance adequada
- **Critérios**:
  - ✅ Consultas paginadas para listas grandes
  - ✅ Índices apropriados no banco de dados
  - ✅ Views otimizadas para consultas complexas
  - ⏳ Tempo de resposta < 2s para operações CRUD
  - ⏳ Suporte a 100+ usuários simultâneos

### RNF02 - Segurança
- **Descrição**: Sistema deve ser seguro
- **Critérios**:
  - ✅ Autenticação via JWT
  - ✅ Autorização baseada em roles
  - ✅ Validação de entrada em todas as APIs
  - ⏳ Hash seguro de senhas (bcrypt)
  - ⏳ Rate limiting para APIs

### RNF03 - Usabilidade
- **Descrição**: Sistema deve ser fácil de usar
- **Critérios**:
  - ✅ API REST bem documentada
  - ✅ Mensagens de erro claras
  - ✅ Códigos de status HTTP apropriados
  - ✅ Formato de resposta consistente
  - ⏳ Interface web responsiva

### RNF04 - Confiabilidade
- **Descrição**: Sistema deve ser confiável
- **Critérios**:
  - ✅ Transações atômicas para operações críticas
  - ✅ Validação de integridade referencial
  - ✅ Rollback automático em caso de erro
  - ⏳ Backup automático do banco
  - ⏳ Logs detalhados de operações

### RNF05 - Manutenibilidade
- **Descrição**: Sistema deve ser fácil de manter
- **Critérios**:
  - ✅ Arquitetura limpa (Clean Architecture)
  - ✅ Separação clara de responsabilidades
  - ✅ Código bem documentado
  - ✅ Testes unitários e de integração
  - ✅ Padrões de código consistentes

### RNF06 - Escalabilidade
- **Descrição**: Sistema deve ser escalável
- **Critérios**:
  - ✅ Arquitetura modular
  - ✅ Separação frontend/backend
  - ✅ APIs stateless
  - ⏳ Suporte a múltiplas instâncias
  - ⏳ Cache para dados frequentes

## 📊 Regras de Negócio Detalhadas

### RN01 - Hierarquia de Entidades
```
Área (1) ──→ (N) Exame
Tópico (N) ←──→ (N) Exame (com configurações)
Tópico (1) ──→ (N) Questão
Exame (N) ←──→ (N) Questão (com validação)
Questão (N) ←──→ (N) Tag
```

### RN02 - Validação Exame-Questão
1. **Regra Principal**: Questão pode ser associada ao exame apenas se:
   - Questão tem tópico atribuído
   - Tópico da questão está associado ao exame

2. **Exceção**: Se exame tem apenas 1 tópico com peso 100%:
   - Qualquer questão pode ser associada
   - Permite flexibilidade para provas gerais

### RN03 - Distribuição de Pesos
- Soma dos pesos de todos os tópicos de um exame = 100%
- Soma das porcentagens de dificuldade por tópico = 100%
- Valores permitidos: 0.01% a 100.00%

### RN04 - Questões e Opções
- Mínimo 2 opções por questão
- Máximo 6 opções por questão (recomendado: 4)
- Exatamente 1 opção correta por questão
- Texto da opção não pode estar vazio

### RN05 - Tags
- Nome da tag deve ser único (case-insensitive)
- Mínimo 2 caracteres por nome
- Máximo 100 caracteres por nome
- Uma questão pode ter 0 a N tags

## 🔄 Casos de Uso Principais

### UC01 - Criar Prova Completa
1. Usuário cria/seleciona área
2. Usuário cria exame na área
3. Usuário associa tópicos ao exame com pesos
4. Usuário cria/seleciona questões para cada tópico
5. Sistema valida e associa questões ao exame
6. Prova está pronta para aplicação

### UC02 - Importar Questões
1. Usuário seleciona arquivo JSON
2. Sistema valida formato e dados
3. Sistema cria/atualiza área, exame, tópico
4. Sistema cria questões e opções
5. Sistema associa tags às questões
6. Sistema valida e associa questões ao exame
7. Importação concluída com sucesso

### UC03 - Reutilizar Questões
1. Usuário cria novo exame
2. Usuário associa tópicos existentes
3. Sistema lista questões disponíveis por tópico
4. Usuário seleciona questões existentes
5. Sistema valida hierarquia tópico-exame
6. Questões são associadas ao novo exame

## 🧪 Cenários de Teste

### CT01 - Validação Hierárquica
```
DADO: Exame com tópicos A e B
QUANDO: Tentar associar questão do tópico C
ENTÃO: Sistema deve rejeitar com erro específico

DADO: Exame com único tópico A (peso 100%)
QUANDO: Tentar associar questão do tópico B
ENTÃO: Sistema deve permitir (exceção)
```

### CT02 - Distribuição de Pesos
```
DADO: Exame com 3 tópicos
QUANDO: Definir pesos 40%, 30%, 35%
ENTÃO: Sistema deve aceitar (soma = 105% > 100%)

QUANDO: Definir pesos 40%, 30%, 30%
ENTÃO: Sistema deve aceitar (soma = 100%)
```

### CT03 - Importação JSON
```
DADO: JSON com área existente e exame novo
QUANDO: Importar arquivo
ENTÃO: Sistema deve usar área existente e criar exame

DADO: JSON com dados inválidos
QUANDO: Importar arquivo
ENTÃO: Sistema deve rejeitar com erro específico
```

## 📈 Métricas de Sucesso

### Funcionalidade
- ✅ 100% dos requisitos funcionais implementados
- ✅ Validação de regras de negócio funcionando
- ⏳ 0 bugs críticos em produção
- ⏳ Tempo médio de resposta < 2s

### Qualidade
- ✅ Cobertura de testes > 80%
- ✅ Código seguindo padrões estabelecidos
- ✅ Documentação completa e atualizada
- ⏳ Disponibilidade > 99%

### Usabilidade
- ⏳ Tempo médio para criar prova < 10min
- ⏳ Taxa de erro de usuário < 5%
- ⏳ Satisfação do usuário > 8/10

## 🔮 Roadmap Futuro

### Versão 2.0
- [ ] Interface web completa
- [ ] Relatórios e estatísticas
- [ ] Aplicação de simulados online
- [ ] Correção automática

### Versão 3.0
- [ ] Integração com LMS
- [ ] Banco de questões público
- [ ] IA para geração de questões
- [ ] Análise de desempenho avançada

---

**Legenda:**
- ✅ Implementado
- ⏳ Planejado
- [ ] Futuro

**Status Atual:** Versão 1.0 - Core completo e funcional
