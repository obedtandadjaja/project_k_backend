Tech debts:
- We can refactor `getTransactionAndQueryContext` into using `pop.Connection.Scope`
- Integrate token endpoint to belong entirely in authService
  - may need to transfer over the userID - does it make sense for credentials to keep userID?
