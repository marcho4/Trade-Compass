ALTER TABLE metrics
    ADD COLUMN company_type VARCHAR(20),
    ADD COLUMN net_interest_income BIGINT,
    ADD COLUMN commission_income BIGINT,
    ADD COLUMN commission_expense BIGINT,
    ADD COLUMN net_commission_income BIGINT,
    ADD COLUMN credit_loss_provision BIGINT;
