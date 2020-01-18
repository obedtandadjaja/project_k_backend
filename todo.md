Tech debts:
- We can refactor `getTransactionAndQueryContext` into using `pop.Connection.Scope`

Features:
- Separate endpoints for tenants: `/api/tenants/{version}/...`
  - create, update, list, show, delete maintenance requests for their room or property
  - account update
  - update users to have type: `admin`, `tenant`. In the future if there is going to be a case where an admin is also a tenant, we can switch to having `types` array instead where we provide the capability for the user to switch between admin and tenant mode.
