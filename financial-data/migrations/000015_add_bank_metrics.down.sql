ALTER TABLE metrics
    DROP COLUMN IF EXISTS company_type,
    DROP COLUMN IF EXISTS net_interest_income,
    DROP COLUMN IF EXISTS commission_income,
    DROP COLUMN IF EXISTS commission_expense,
    DROP COLUMN IF EXISTS net_commission_income,
    DROP COLUMN IF EXISTS credit_loss_provision;
