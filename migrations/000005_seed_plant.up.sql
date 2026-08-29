INSERT INTO plants (epoch, phase, phase_progress, branching, density, curvature, vitality, seed)
SELECT 0, 0, 0, 48, 71, 22, 91, 12345
WHERE NOT EXISTS (SELECT 1 FROM plants);
