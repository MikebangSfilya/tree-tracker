DELETE FROM plants WHERE id = (SELECT MAX(id) FROM plants);
