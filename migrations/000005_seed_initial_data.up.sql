INSERT INTO users (
    id,
    username,
    email,
    password
)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'system',
    'system@tree-tracker.local',
    convert_to('$2a$10$JvKTAYglrzCoaTMEQwHPJetkwuns.18MDs8ML2pkIWvvYCkXtCSQu', 'UTF8')
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO plants (
    epoch,
    phase,
    phase_progress,
    branching,
    density,
    curvature,
    vitality,
    seed
)
VALUES (
    0,
    1,
    37,
    48,
    31,
    22,
    94,
    12345
);
