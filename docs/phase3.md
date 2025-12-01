
```markdown
# Phase 3: Product Management System

## Overview

Built a complete product catalog system with categories, tags, and products. Includes advanced filtering, search, inventory management, and relationship handling.

## What We Built

### 1. Category Model (`internal/models/category.go`)

**Hierarchical category system with parent-child relationships.**

- **GORM Model** with proper struct tags
- Fields: ID, Name, Slug, Description, ParentID, ImageURL, IsActive, Timestamps
- **Parent-Child Support** - Nested categories (unlimited depth)
- **Soft Delete** support with `DeletedAt`
- **ToResponse()** method for safe API responses
- **ToResponseWithChildren()** for hierarchical data

**Key Features:**
- Slug-based URLs for SEO (`/categories/electronics`)
- Self-referencing foreign key for hierarchy
- Automatic timestamp management
- Prevents deletion if has children
- Tree structure support

**Example Structure:**
```
Electronics (parent_id: NULL)
├── Phones (parent_id: 1)
│   ├── Smartphones (parent_id: 2)
│   └── Feature Phones (parent_id: 2)
└── Laptops (parent_id: 1)
    ├── Gaming Laptops (parent_id: 5)
    └── Business Laptops (parent_id: 5)
```

---

### 2. Tag Model (`internal/models/tag.go`)

**Simple labeling system for products.**

- **GORM Model** with proper struct tags
- Fields: ID, Name, Slug, Timestamps
- **Many-to-Many** relationship with products
- **Soft Delete** support
- **ToResponse()** method for API responses

**Key Features:**
- Unique slug per tag
- Can be shared across multiple products
- Search functionality
- Lightweight and fast

**Example Tags:**
```
- new-arrival
- on-sale
- trending
- best-seller
- limited-edition
```

---

### 3. Product Model (`internal/models/product.go`)

**Complete product entity with rich features.**

- **GORM Model** with proper struct tags
- Fields: ID, Name, Slug, Description, Price, ComparePrice, SKU, Barcode, Quantity, Images, CategoryID, IsActive, IsFeatured, Timestamps
- **Belongs To** Category (foreign key)
- **Many-to-Many** with Tags
- **JSON Images** - Array of image URLs
- **Soft Delete** support
- **ToResponse()** with calculated discount percentage

**Product Fields:**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Product name |
| `slug` | string | URL-friendly identifier |
| `description` | text | Product details |
| `price` | decimal(10,2) | Current selling price |
| `compare_price` | decimal(10,2) | Original price (for showing discount) |
| `sku` | string | Stock Keeping Unit (unique) |
| `barcode` | string | Product barcode |
| `quantity` | int | Available stock |
| `images` | JSON | Array of image URLs |
| `category_id` | uint | Foreign key to category |
| `is_active` | bool | Product visibility |
| `is_featured` | bool | Show in featured section |

**Calculated Fields:**
- `discount` - Percentage off (calculated from price vs compare_price)
- `in_stock` - Boolean (quantity > 0)

**Example Product:**
```
{
  "id": 1,
  "name": "iPhone 15 Pro Max",
  "slug": "iphone-15-pro-max",
  "price": 1199.00,
  "compare_price": 1299.00,
  "discount": 7.7,
  "sku": "APPL-IPH15-PM-256-BLK",
  "quantity": 50,
  "in_stock": true,
  "images": [
    "https://cdn.example.com/iphone15-1.jpg",
    "https://cdn.example.com/iphone15-2.jpg"
  ],
  "category": {
    "id": 2,
    "name": "Smartphones",
    "slug": "smartphones"
  },
  "tags": [
    {"id": 1, "name": "New Arrival", "slug": "new-arrival"},
    {"id": 3, "name": "Trending", "slug": "trending"}
  ]
}
```

---

### 4. Category Repository (`internal/repository/category.go`)

**Database access layer for categories.**

**Methods:**

| Method | Description |
|--------|-------------|
| `Create(category)` | Insert new category |
| `GetByID(id)` | Find by ID |
| `GetBySlug(slug)` | Find by slug |
| `GetAll(limit, offset)` | Paginated list |
| `GetRootCategories()` | Top-level categories only |
| `GetWithChildren(id)` | Category with subcategories |
| `GetCategoryTree()` | Full hierarchy |
| `GetActiveCategories()` | Only active categories |
| `Update(category)` | Update category |
| `Delete(id)` | Soft delete |
| `SlugExists(slug)` | Check uniqueness |

**Special Queries:**
- Tree structure with `Preload("Children")`
- Root categories: `WHERE parent_id IS NULL`
- Active filtering: `WHERE is_active = true`

---

### 5. Tag Repository (`internal/repository/tag.go`)

**Database access layer for tags.**

**Methods:**

| Method | Description |
|--------|-------------|
| `Create(tag)` | Insert new tag |
| `GetByID(id)` | Find by ID |
| `GetBySlug(slug)` | Find by slug |
| `GetAll(limit, offset)` | Paginated list |
| `GetAllTags()` | All tags (no pagination) |
| `Search(query, limit, offset)` | Search by name |
| `Update(tag)` | Update tag |
| `Delete(id)` | Soft delete |
| `SlugExists(slug)` | Check uniqueness |

**Search Example:**
```
// Find tags containing "new"
tags, total, _ := tagRepo.Search("new", 10, 0)
// Returns: ["New Arrival", "Renewed", "New Collection"]
```

---

### 6. Product Repository (`internal/repository/product.go`)

**Database access layer for products with advanced filtering.**

**CRUD Methods:**

| Method | Description |
|--------|-------------|
| `Create(product)` | Insert new product |
| `GetByID(id)` | Find by ID (with category & tags) |
| `GetBySlug(slug)` | Find by slug (with relations) |
| `GetBySKU(sku)` | Find by SKU |
| `GetAll(limit, offset, filters)` | Advanced filtered list |
| `Update(product)` | Update product |
| `UpdateTags(product, tags)` | Update many-to-many tags |
| `Delete(id)` | Soft delete |

**Special Methods:**

| Method | Description |
|--------|-------------|
| `GetFeaturedProducts(limit)` | Featured products only |
| `GetRelatedProducts(id, categoryID, limit)` | Same category products |
| `UpdateStock(id, quantity)` | Set stock quantity |
| `DecrementStock(id, amount)` | Reduce stock (for orders) |
| `SlugExists(slug)` | Check slug uniqueness |
| `SKUExists(sku)` | Check SKU uniqueness |

**Advanced Filtering:**

```
type ProductFilterOptions struct {
    CategoryID  *uint      // Filter by category
    TagIDs      []uint     // Filter by tags (OR)
    MinPrice    *float64   // Price range min
    MaxPrice    *float64   // Price range max
    InStock     *bool      // Only in-stock items
    IsActive    *bool      // Only active products
    IsFeatured  *bool      // Only featured
    SearchQuery string     // Search name/description/SKU
}
```

**Filter Examples:**
```
// Products in category 5, price $100-$500, in stock
filters := &ProductFilterOptions{
    CategoryID: &categoryID,  // 5
    MinPrice:   &minPrice,    // 100.0
    MaxPrice:   &maxPrice,    // 500.0
    InStock:    &inStock,     // true
}

// Search for "phone" with tags 
filters := &ProductFilterOptions{
    SearchQuery: "phone",
    TagIDs:      []uint{1, 3},
}
```

---

### 7. Category Service (`internal/service/category.go`)

**Business logic for categories.**

**DTOs:**
- `CreateCategoryRequest` - Validation rules
- `UpdateCategoryRequest` - Optional fields
- `CategoryListResponse` - Pagination metadata

**Service Methods:**

| Method | Business Logic |
|--------|----------------|
| `CreateCategory()` | Generate slug, check uniqueness, validate parent |
| `GetCategoryByID()` | Retrieve single category |
| `GetCategoryBySlug()` | SEO-friendly retrieval |
| `GetCategoryWithChildren()` | Category + subcategories |
| `GetAllCategories()` | Paginated list with metadata |
| `GetActiveCategories()` | Only active categories |
| `GetRootCategories()` | Top-level categories |
| `GetCategoryTree()` | Full hierarchy |
| `UpdateCategory()` | Regenerate slug, validate parent, prevent self-reference |
| `DeleteCategory()` | Check for children, prevent deletion if has subcategories |

**Validation Rules:**
- ✅ Name: 2-255 characters
- ✅ Slug: Auto-generated, unique
- ✅ Parent: Must exist, cannot self-reference
- ✅ Deletion: Blocked if has children

---

### 8. Tag Service (`internal/service/tag.go`)

**Business logic for tags.**

**DTOs:**
- `CreateTagRequest` - Name validation
- `UpdateTagRequest` - Name update
- `TagListResponse` - Pagination metadata

**Service Methods:**

| Method | Business Logic |
|--------|----------------|
| `CreateTag()` | Generate slug, check uniqueness |
| `GetTagByID()` | Retrieve single tag |
| `GetTagBySlug()` | SEO-friendly retrieval |
| `GetAllTags()` | Paginated list |
| `GetAllTagsList()` | All tags (for dropdowns) |
| `SearchTags()` | Search by name |
| `UpdateTag()` | Regenerate slug, check uniqueness |
| `DeleteTag()` | Soft delete (TODO: check product usage) |

**Validation Rules:**
- ✅ Name: 2-100 characters
- ✅ Slug: Auto-generated, unique

---

### 9. Product Service (`internal/service/product.go`)

**Business logic for products with complex validation.**

**DTOs:**
- `CreateProductRequest` - Full validation
- `UpdateProductRequest` - All fields optional
- `ProductFilterRequest` - Query parameters
- `ProductListResponse` - Pagination metadata

**Service Methods:**

| Method | Business Logic |
|--------|----------------|
| `CreateProduct()` | Validate category, tags, generate slug, check SKU uniqueness, convert images to JSON |
| `GetProductByID()` | With category and tags preloaded |
| `GetProductBySlug()` | SEO-friendly retrieval |
| `GetAllProducts()` | Advanced filtering, pagination |
| `GetFeaturedProducts()` | Limited featured list |
| `GetRelatedProducts()` | Same category, random selection |
| `UpdateProduct()` | Update fields, regenerate slug, validate SKU, update tags association |
| `DeleteProduct()` | Soft delete (TODO: check orders) |
| `UpdateStock()` | Inventory management |

**Validation Rules:**

| Field | Validation |
|-------|-----------|
| `name` | Required, 2-255 chars |
| `price` | Required, > 0 |
| `compare_price` | Optional, >= price |
| `sku` | Optional, unique, max 100 |
| `quantity` | Optional, >= 0 |
| `category_id` | Required, must exist |
| `tag_ids` | Optional, all must exist |
| `images` | Optional, array of URLs |

**JSON Image Handling:**
```
// Convert []string to JSON for database
images := []string{
    "https://cdn.example.com/img1.jpg",
    "https://cdn.example.com/img2.jpg"
}
imagesBytes, _ := json.Marshal(images)
product.Images = datatypes.JSON(imagesBytes)
```

---

### 10. Slug Utility (`internal/utils/slug.go`)

**URL-friendly slug generation.**

**Function:** `GenerateSlug(s string)`

**Features:**
- Converts to lowercase
- Removes accents/diacritics
- Replaces spaces with hyphens
- Removes special characters
- Collapses multiple hyphens

**Examples:**
```
GenerateSlug("iPhone 15 Pro Max")       → "iphone-15-pro-max"
GenerateSlug("Men's Fashion")           → "mens-fashion"
GenerateSlug("Électronique & Gadgets")  → "electronique-gadgets"
GenerateSlug("Gaming   Laptop!!!")      → "gaming-laptop"
```

---

## Database Schema

### **categories table:**
```
CREATE TABLE categories (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  parent_id BIGINT UNSIGNED,
  image_url VARCHAR(500),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_categories_slug (slug),
  INDEX idx_categories_parent (parent_id),
  INDEX idx_categories_deleted_at (deleted_at),
  FOREIGN KEY (parent_id) REFERENCES categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### **tags table:**
```
CREATE TABLE tags (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(100) NOT NULL,
  slug VARCHAR(100) NOT NULL UNIQUE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_tags_slug (slug),
  INDEX idx_tags_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### **products table:**
```
CREATE TABLE products (
  id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  description TEXT,
  price DECIMAL(10,2) NOT NULL,
  compare_price DECIMAL(10,2),
  sku VARCHAR(100) UNIQUE,
  barcode VARCHAR(100),
  quantity INT DEFAULT 0,
  images JSON,
  category_id BIGINT UNSIGNED,
  is_active BOOLEAN DEFAULT TRUE,
  is_featured BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP,
  updated_at TIMESTAMP,
  deleted_at TIMESTAMP NULL,
  
  INDEX idx_products_slug (slug),
  INDEX idx_products_sku (sku),
  INDEX idx_products_category (category_id),
  INDEX idx_products_deleted_at (deleted_at),
  FOREIGN KEY (category_id) REFERENCES categories(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### **product_tags table (many-to-many):**
```
CREATE TABLE product_tags (
  product_id BIGINT UNSIGNED,
  tag_id BIGINT UNSIGNED,
  PRIMARY KEY (product_id, tag_id),
  
  FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## API Endpoints

### **Public Endpoints (No Authentication)**

#### **Categories:**
```
GET    /api/v1/categories                    → List all (paginated)
GET    /api/v1/categories/active             → Active only
GET    /api/v1/categories/root               → Top-level only
GET    /api/v1/categories/tree               → Full hierarchy
GET    /api/v1/categories/:id                → Get by ID
GET    /api/v1/categories/slug/:slug         → Get by slug
GET    /api/v1/categories/:id/children       → With subcategories
```

#### **Tags:**
```
GET    /api/v1/tags                          → List all (paginated)
GET    /api/v1/tags/all                      → All tags (no pagination)
GET    /api/v1/tags/search?q=keyword         → Search tags
GET    /api/v1/tags/:id                      → Get by ID
GET    /api/v1/tags/slug/:slug               → Get by slug
```

#### **Products:**
```
GET    /api/v1/products                      → List all (with filters)
GET    /api/v1/products/featured             → Featured products
GET    /api/v1/products/:id                  → Get by ID
GET    /api/v1/products/slug/:slug           → Get by slug
GET    /api/v1/products/:id/related          → Related products
```

**Product Filter Parameters:**
```
?page=1
?limit=20
?category_id=5
?tag_ids=1,3,5
?min_price=100
?max_price=500
?in_stock=true
?is_active=true
?is_featured=true
?q=search+query
```

---

### **Admin Endpoints (Requires JWT + Admin Role)**

#### **Categories:**
```
POST   /api/v1/admin/categories              → Create category
PUT    /api/v1/admin/categories/:id          → Update category
DELETE /api/v1/admin/categories/:id          → Delete category
```

#### **Tags:**
```
POST   /api/v1/admin/tags                    → Create tag
PUT    /api/v1/admin/tags/:id                → Update tag
DELETE /api/v1/admin/tags/:id                → Delete tag
```

#### **Products:**
```
POST   /api/v1/admin/products                → Create product
PUT    /api/v1/admin/products/:id            → Update product
DELETE /api/v1/admin/products/:id            → Delete product
PUT    /api/v1/admin/products/:id/stock      → Update stock
```

---

## Complete API Examples

### **1. Create Category (Admin)**

**Request:**
```
curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and gadgets",
    "is_active": true
  }'
```

**Response (201):**
```
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": 1,
    "name": "Electronics",
    "slug": "electronics",
    "description": "Electronic devices and gadgets",
    "is_active": true,
    "created_at": "2025-12-02T00:30:00Z",
    "updated_at": "2025-12-02T00:30:00Z"
  }
}
```

---

### **2. Create Subcategory (Admin)**

**Request:**
```
curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Smartphones",
    "description": "Mobile phones and accessories",
    "parent_id": 1,
    "is_active": true
  }'
```

---

### **3. Create Tag (Admin)**

**Request:**
```
curl -X POST http://localhost:8080/api/v1/admin/tags \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "New Arrival"
  }'
```

**Response (201):**
```
{
  "success": true,
  "message": "Tag created successfully",
  "data": {
    "id": 1,
    "name": "New Arrival",
    "slug": "new-arrival",
    "created_at": "2025-12-02T00:31:00Z",
    "updated_at": "2025-12-02T00:31:00Z"
  }
}
```

---

### **4. Create Product (Admin)**

**Request:**
```
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro Max",
    "description": "Latest flagship smartphone from Apple",
    "price": 1199.00,
    "compare_price": 1299.00,
    "sku": "APPL-IPH15-PM-256",
    "quantity": 50,
    "images": [
      "https://cdn.example.com/iphone15-1.jpg",
      "https://cdn.example.com/iphone15-2.jpg",
      "https://cdn.example.com/iphone15-3.jpg"
    ],
    "category_id": 2,
    "tag_ids": ,
    "is_active": true,
    "is_featured": true
  }'
```

**Response (201):**
```
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": 1,
    "name": "iPhone 15 Pro Max",
    "slug": "iphone-15-pro-max",
    "description": "Latest flagship smartphone from Apple",
    "price": 1199.00,
    "compare_price": 1299.00,
    "discount": 7.7,
    "sku": "APPL-IPH15-PM-256",
    "quantity": 50,
    "in_stock": true,
    "images": [
      "https://cdn.example.com/iphone15-1.jpg",
      "https://cdn.example.com/iphone15-2.jpg",
      "https://cdn.example.com/iphone15-3.jpg"
    ],
    "category": {
      "id": 2,
      "name": "Smartphones",
      "slug": "smartphones"
    },
    "tags": [
      {"id": 1, "name": "New Arrival", "slug": "new-arrival"},
      {"id": 3, "name": "Trending", "slug": "trending"}
    ],
    "is_active": true,
    "is_featured": true,
    "created_at": "2025-12-02T00:32:00Z",
    "updated_at": "2025-12-02T00:32:00Z"
  }
}
```

---

### **5. Get All Products (Public)**

**Request:**
```
curl "http://localhost:8080/api/v1/products?page=1&limit=10"
```

**Response (200):**
```
{
  "success": true,
  "message": "Products retrieved successfully",
  "data": {
    "products": [
      {
        "id": 1,
        "name": "iPhone 15 Pro Max",
        "slug": "iphone-15-pro-max",
        "price": 1199.00,
        "discount": 7.7,
        "in_stock": true,
        "images": ["https://cdn.example.com/iphone15-1.jpg"],
        "category": {"id": 2, "name": "Smartphones"},
        "tags": [{"id": 1, "name": "New Arrival"}]
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 10,
    "total_pages": 1
  }
}
```

---

### **6. Filter Products (Public)**

**Request:**
```
# Products in category 2, price $500-$1500, in stock, with tag "new-arrival"
curl "http://localhost:8080/api/v1/products?category_id=2&min_price=500&max_price=1500&in_stock=true&tag_ids=1"
```

---

### **7. Search Products (Public)**

**Request:**
```
curl "http://localhost:8080/api/v1/products?q=iphone"
```

---

### **8. Get Category Tree (Public)**

**Request:**
```
curl http://localhost:8080/api/v1/categories/tree
```

**Response (200):**
```
{
  "success": true,
  "message": "Category tree retrieved successfully",
  "data": [
    {
      "id": 1,
      "name": "Electronics",
      "slug": "electronics",
      "children": [
        {
          "id": 2,
          "name": "Smartphones",
          "slug": "smartphones"
        },
        {
          "id": 3,
          "name": "Laptops",
          "slug": "laptops"
        }
      ]
    }
  ]
}
```

---

### **9. Get Featured Products (Public)**

**Request:**
```
curl "http://localhost:8080/api/v1/products/featured?limit=5"
```

---

### **10. Get Related Products (Public)**

**Request:**
```
curl "http://localhost:8080/api/v1/products/1/related?limit=4"
```

---

### **11. Update Product Stock (Admin)**

**Request:**
```
curl -X PUT http://localhost:8080/api/v1/admin/products/1/stock \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 100
  }'
```

**Response (200):**
```
{
  "success": true,
  "message": "Stock updated successfully",
  "data": {
    "product_id": 1,
    "quantity": 100
  }
}
```

---

## Architecture Pattern

**Clean Architecture / Layered Architecture:**

```
Handler (HTTP Layer)
    ↓
Service (Business Logic)
    ↓
Repository (Database Layer)
    ↓
Database (GORM/MySQL)
```

**Benefits:**
- ✅ Separation of Concerns
- ✅ Testable (each layer can be mocked)
- ✅ Maintainable (easy to modify one layer)
- ✅ Scalable (can swap implementations)

**Dependency Injection:**
```
// main.go
productRepo := repository.NewProductRepository(db)
categoryRepo := repository.NewCategoryRepository(db)
tagRepo := repository.NewTagRepository(db)

productService := service.NewProductService(productRepo, categoryRepo, tagRepo)
productHandler := handler.NewProductHandler(productService)
```

---

## Key Features Implemented

### **1. Hierarchical Categories**
- ✅ Unlimited nesting depth
- ✅ Parent-child relationships
- ✅ Tree structure queries
- ✅ Prevents circular references
- ✅ Cascade protection (can't delete with children)

### **2. Product Tagging System**
- ✅ Many-to-many relationships
- ✅ Multiple tags per product
- ✅ Tag-based filtering
- ✅ Search functionality

### **3. Advanced Product Filtering**
- ✅ Filter by category
- ✅ Filter by tags (OR logic)
- ✅ Price range filtering
- ✅ Stock availability filter
- ✅ Active/Featured filters
- ✅ Full-text search (name, description, SKU)
- ✅ Combine multiple filters

### **4. SEO-Friendly URLs**
- ✅ Slug generation from names
- ✅ URL-safe characters only
- ✅ Unique slugs enforced
- ✅ Access by slug: `/products/slug/iphone-15-pro-max`

### **5. Image Management**
- ✅ Multiple images per product
- ✅ JSON array storage
- ✅ Easy to add/remove images
- ✅ Image URL validation

### **6. Inventory Tracking**
- ✅ Stock quantity management
- ✅ In-stock calculation
- ✅ Stock update endpoint
- ✅ Decrement stock (prepared for orders)

### **7. Pricing & Discounts**
- ✅ Current price
- ✅ Compare price (original price)
- ✅ Auto-calculated discount percentage
- ✅ Price range filtering

### **8. Product Discovery**
- ✅ Featured products
- ✅ Related products (same category)
- ✅ Search functionality
- ✅ Category browsing

### **9. Soft Delete**
- ✅ Categories can be recovered
- ✅ Tags can be recovered
- ✅ Products can be recovered
- ✅ Excluded from queries by default

### **10. Pagination**
- ✅ All list endpoints support pagination
- ✅ Total count metadata
- ✅ Page calculation
- ✅ Configurable limits

---

## Files Created in Phase 3

```
internal/models/
  ├── category.go              ✅ Hierarchical category model
  ├── tag.go                   ✅ Simple tag model
  └── product.go               ✅ Complex product model with relations

internal/repository/
  ├── category.go              ✅ Category CRUD + tree queries
  ├── tag.go                   ✅ Tag CRUD + search
  └── product.go               ✅ Product CRUD + advanced filters

internal/service/
  ├── category.go              ✅ Category business logic
  ├── tag.go                   ✅ Tag business logic
  └── product.go               ✅ Product business logic + validation

internal/handler/
  ├── category.go              ✅ Category HTTP handlers
  ├── tag.go                   ✅ Tag HTTP handlers
  └── product.go               ✅ Product HTTP handlers

internal/utils/
  └── slug.go                  ✅ URL slug generation

cmd/api/
  └── main.go                  ✅ Updated with all routes

internal/database/
  └── migrate.go               ✅ Updated with Product model
```

---

## Testing Guide

### **Setup Test Data**

#### **1. Create Admin User:**
```
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Admin",
    "last_name": "User",
    "email": "admin@shop.com",
    "password": "admin123456"
  }'

# Verify (check logs for token)
curl "http://localhost:8080/api/v1/auth/verify?token=TOKEN_FROM_LOGS"

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@shop.com",
    "password": "admin123456"
  }'
```

**Set user as admin in database:**
```
USE ecommerce;
UPDATE users SET role = 'admin' WHERE email = 'admin@shop.com';
```

Save the JWT token for admin requests!

---

#### **2. Create Category Hierarchy:**

```
TOKEN="your_admin_jwt_token"

# Root category
curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Electronics"}'

# Subcategories (assuming Electronics has id=1)
curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Smartphones", "parent_id": 1}'

curl -X POST http://localhost:8080/api/v1/admin/categories \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Laptops", "parent_id": 1}'
```

---

#### **3. Create Tags:**

```
curl -X POST http://localhost:8080/api/v1/admin/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "New Arrival"}'

curl -X POST http://localhost:8080/api/v1/admin/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "On Sale"}'

curl -X POST http://localhost:8080/api/v1/admin/tags \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Trending"}'
```

---

#### **4. Create Products:**

```
# iPhone
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15 Pro Max",
    "description": "Latest flagship smartphone",
    "price": 1199.00,
    "compare_price": 1299.00,
    "sku": "IPH15PM-256",
    "quantity": 50,
    "category_id": 2,
    "tag_ids": ,
    "is_featured": true
  }'

# MacBook
curl -X POST http://localhost:8080/api/v1/admin/products \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "MacBook Pro 16 M3",
    "description": "Powerful laptop for professionals",
    "price": 2499.00,
    "sku": "MBP16-M3-512",
    "quantity": 30,
    "category_id": 3,
    "tag_ids": ,
    "is_featured": true
  }'
```

---

### **Test Cases**

#### **Categories:**
- [ ] Create root category
- [ ] Create subcategory
- [ ] Get category tree
- [ ] Get category by slug
- [ ] Update category name (slug regenerates)
- [ ] Try to delete category with children (should fail)
- [ ] Delete leaf category (should succeed)

#### **Tags:**
- [ ] Create tag
- [ ] Search tags by name
- [ ] Get all tags (no pagination)
- [ ] Update tag name
- [ ] Delete tag

#### **Products:**
- [ ] Create product with category
- [ ] Create product with tags
- [ ] Create product with images
- [ ] Get product by ID (includes category & tags)
- [ ] Get product by slug
- [ ] Filter by category
- [ ] Filter by price range
- [ ] Filter by tags
- [ ] Search by keyword
- [ ] Get featured products
- [ ] Get related products
- [ ] Update product
- [ ] Update stock
- [ ] Delete product

---

## Common Issues & Solutions

**Issue:** "category not found" when creating product
- **Solution:** Create category first, use correct category_id

**Issue:** "tag with ID X not found"
- **Solution:** Create tags first, use correct tag IDs

**Issue:** "SKU already exists"
- **Solution:** Use unique SKU for each product

**Issue:** "category with similar name already exists"
- **Solution:** Slugs must be unique, change category name

**Issue:** Cannot delete category
- **Solution:** Delete child categories first, or reassign them

**Issue:** Products not appearing in filters
- **Solution:** Check is_active=true, quantity>0 for in_stock filter

---

## Environment Variables

No new variables needed for Phase 3! Uses existing database config.

```
# Database (from Phase 1)
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=ecommerce

# JWT (from Phase 2)
JWT_SECRET=your-secret-key
JWT_EXPIRATION=24h
```

---

## Performance Considerations

### **Eager Loading:**
```
// Load category and tags with product (avoids N+1 queries)
db.Preload("Category").Preload("Tags").Find(&products)
```

### **Indexes:**
- ✅ Slug columns (unique index)
- ✅ Category foreign keys
- ✅ SKU (unique index)
- ✅ Deleted_at (for soft delete queries)

### **Pagination:**
- ✅ Always use LIMIT/OFFSET
- ✅ Return total count for UI
- ✅ Default limits prevent large queries

### **Query Optimization:**
- ✅ Filter before loading relationships
- ✅ Use specific field selection when needed
- ✅ Index foreign keys

---

## Future Enhancements (Phase 4+)

- [ ] Product variants (size, color, etc.)
- [ ] Product reviews & ratings
- [ ] Image upload/storage (S3, Cloudinary)
- [ ] Advanced search (Elasticsearch)
- [ ] Product recommendations (AI)
- [ ] Inventory alerts (low stock notifications)
- [ ] Bulk import/export (CSV)
- [ ] Product analytics
- [ ] Wishlist functionality
- [ ] Compare products

---

## Security Features

### **Admin Authorization:**
```
// All admin endpoints check role
userRole, exists := middleware.GetUserRole(c)
if !exists || userRole != "admin" {
    utils.RespondForbidden(c, "Admin access required")
    return
}
```

### **Input Validation:**
- ✅ Gin binding tags (`required`, `min`, `max`, `gt`, `gte`)
- ✅ Price validation (must be > 0)
- ✅ Quantity validation (must be >= 0)
- ✅ Foreign key validation (category/tags exist)

### **SQL Injection Prevention:**
- ✅ GORM parameterized queries
- ✅ No raw SQL with user input

### **Data Sanitization:**
- ✅ Slug generation removes special chars
- ✅ JSON marshaling for images

---

## Key Learnings

1. **Hierarchical Data** - Self-referencing foreign keys for tree structures
2. **Many-to-Many** - Join tables for flexible relationships
3. **JSON Storage** - Storing arrays in relational databases
4. **Eager Loading** - Preload to avoid N+1 query problem
5. **Slug Generation** - SEO-friendly URL patterns
6. **Advanced Filtering** - Building complex WHERE clauses
7. **Discount Calculation** - Computed fields in responses
8. **Clean Architecture** - Separation of concerns across layers

---

## Phase 3: Complete! ✅

**What We Built:**
- ✅ Hierarchical category system
- ✅ Flexible tagging system
- ✅ Full-featured product catalog
- ✅ Advanced filtering & search
- ✅ Inventory management
- ✅ SEO-friendly URLs
- ✅ Admin CRUD operations
- ✅ Public browsing APIs
- ✅ Relationship management
- ✅ Clean architecture

**Database Tables:**
- ✅ `categories` (hierarchical)
- ✅ `tags` (simple)
- ✅ `products` (complex)
- ✅ `product_tags` (many-to-many join)

**Ready for Phase 4: Shopping Cart & Orders!** 🛒🚀
```