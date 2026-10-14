CREATE TABLE criteria_configs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    config TEXT NOT NULL CHECK (json_valid(config)),
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

INSERT INTO criteria_configs (name, config, created_at, updated_at)
SELECT
    'Default',
    '{
        "common": {
            "label": "Style Criteria",
            "tags": [
                {
                    "id": "common_initiation",
                    "label": "Initiation",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 5
                },
                {
                    "id": "common_fluidity",
                    "label": "Fluidity",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 10
                },
                {
                    "id": "common_commitment",
                    "label": "Commitment",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 5
                }
            ]
        },
        "lead": {
            "label": "Lead Criteria",
            "tags": [
                {
                    "id": "lead_filled_n_zones",
                    "label": "Filled Zones",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 30,
                    "multiplier": 7.5
                },
                {
                    "id": "lead_angle",
                    "label": "Angle",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 20
                },
                {
                    "id": "lead_line_quality",
                    "label": "Line Quality",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 30,
                    "multiplier": 6
                },
                {
                    "id": "lead_wheels_off",
                    "label": "Wheels Off",
                    "type": "deduction",
                    "valueType": "dynamic",
                    "max": 100,
                    "multiplier": 25
                },
                {
                    "id": "lead_dnf",
                    "label": "Do Not Finish",
                    "type": "deduction",
                    "valueType": "switch",
                    "defaultValue": -999
                }
            ]
        },
        "chase": {
            "label": "Chase Criteria",
            "tags": [
                {
                    "id": "chase_proximity_quality",
                    "label": "Proximity Quality",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 10
                },
                {
                    "id": "chase_angle",
                    "label": "Angle",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 20
                },
                {
                    "id": "chase_line_quality",
                    "label": "Line Quality",
                    "type": "neutral",
                    "valueType": "dynamic",
                    "max": 10
                },
                {
                    "id": "chase_wheels_off",
                    "label": "Wheels Off",
                    "type": "deduction",
                    "valueType": "dynamic",
                    "max": 100,
                    "multiplier": 25
                },
                {
                    "id": "chase_dnf",
                    "label": "Do Not Finish",
                    "type": "deduction",
                    "valueType": "switch",
                    "defaultValue": -999
                }
            ]
        }
    }',
    unixepoch(),
    unixepoch()
WHERE NOT EXISTS (
    SELECT 1 FROM criteria_configs
);