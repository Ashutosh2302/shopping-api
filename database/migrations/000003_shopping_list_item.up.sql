CREATE TABLE Shopping_List_Item (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    shoppingListId UUID NOT NULL,
    name VARCHAR(50) NOT NULL,
    picked BOOLEAN DEFAULT FALSE,
    price DECIMAL(15, 2) DEFAULT 0.0,
    createdAt TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(shoppingListId, name),
	PRIMARY KEY (id),
    FOREIGN KEY (shoppingListId) REFERENCES Shopping_List(id) ON DELETE CASCADE
);