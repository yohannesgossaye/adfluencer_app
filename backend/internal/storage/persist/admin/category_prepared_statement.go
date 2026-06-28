package admin

const (
	listCategories = `
	SELECT
		c.id,
		c.name,
		c.slug,
		c.description,
		c.is_active,
		c.created_at,
		c.updated_at,
		(SELECT COUNT(DISTINCT ipc.influencer_user_id)::int
		 FROM influencer_profile_categories ipc
		 WHERE ipc.category_id = c.id),
		(SELECT COUNT(*)::int
		 FROM campaigns camp
		 WHERE camp.category_id = c.id)
	FROM categories c
	ORDER BY c.name ASC
	`

	getCategoryByID = `
	SELECT
		c.id,
		c.name,
		c.slug,
		c.description,
		c.is_active,
		c.created_at,
		c.updated_at,
		(SELECT COUNT(DISTINCT ipc.influencer_user_id)::int
		 FROM influencer_profile_categories ipc
		 WHERE ipc.category_id = c.id),
		(SELECT COUNT(*)::int
		 FROM campaigns camp
		 WHERE camp.category_id = c.id)
	FROM categories c
	WHERE c.id = $1
	`

	insertCategory = `
	INSERT INTO categories (id, name, slug, description, is_active)
	VALUES ($1, $2, $3, $4, TRUE)
	RETURNING id, name, slug, description, is_active, created_at, updated_at
	`

	updateCategory = `
	WITH updated AS (
		UPDATE categories
		SET
			name = COALESCE($2, name),
			slug = COALESCE($3, slug),
			description = COALESCE($4, description),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id
	)
	SELECT
		c.id,
		c.name,
		c.slug,
		c.description,
		c.is_active,
		c.created_at,
		c.updated_at,
		(SELECT COUNT(DISTINCT ipc.influencer_user_id)::int
		 FROM influencer_profile_categories ipc
		 WHERE ipc.category_id = c.id),
		(SELECT COUNT(*)::int
		 FROM campaigns camp
		 WHERE camp.category_id = c.id)
	FROM categories c
	INNER JOIN updated u ON c.id = u.id
	`

	toggleCategoryActive = `
	UPDATE categories
	SET is_active = $2, updated_at = NOW()
	WHERE id = $1
	`

	categorySlugExists = `
	SELECT EXISTS (
		SELECT 1 FROM categories WHERE LOWER(slug) = LOWER($1) AND ($2 = '' OR id::text <> $2)
	)
	`
)
