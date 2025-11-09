-- Create plan_features table
CREATE TABLE IF NOT EXISTS plan_features (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    feature_id UUID NOT NULL REFERENCES features(id) ON DELETE CASCADE,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(plan_id, feature_id)
);

-- Create index on plan_id
CREATE INDEX IF NOT EXISTS idx_plan_features_plan_id ON plan_features(plan_id);

-- Create index on feature_id
CREATE INDEX IF NOT EXISTS idx_plan_features_feature_id ON plan_features(feature_id);

-- Create index on enabled
CREATE INDEX IF NOT EXISTS idx_plan_features_enabled ON plan_features(enabled);

-- Create index on deleted_at for soft delete
CREATE INDEX IF NOT EXISTS idx_plan_features_deleted_at ON plan_features(deleted_at);

-- Create trigger to update updated_at column
CREATE OR REPLACE FUNCTION update_plan_features_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER plan_features_updated_at
    BEFORE UPDATE ON plan_features
    FOR EACH ROW
    EXECUTE FUNCTION update_plan_features_updated_at();