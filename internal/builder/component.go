package builder

import (
	"fmt"
	"strings"
)

// Component represents a UI component in the visual builder
type Component struct {
	ID       string            `json:"id"`
	Type     string            `json:"type"`
	Content  string            `json:"content"`
	Styles   map[string]string `json:"styles,omitempty"`
	Children []Component       `json:"children,omitempty"`
}

// ToHTML renders the component as an HTML string
func (c *Component) ToHTML() string {
	styleAttr := buildStyleAttr(c.Styles)
	classes := c.Styles["class"]

	switch c.Type {
	case "heading":
		text := c.Content
		if text == "" {
			text = "Heading"
		}
		return fmt.Sprintf(`<h1 class="text-3xl font-bold %s" style="%s">%s</h1>`, classes, styleAttr, text)

	case "paragraph":
		text := c.Content
		if text == "" {
			text = "Paragraph text here..."
		}
		return fmt.Sprintf(`<p class="text-base-content %s" style="%s">%s</p>`, classes, styleAttr, text)

	case "button":
		text := c.Content
		if text == "" {
			text = "Click Me"
		}
		return fmt.Sprintf(`<button class="btn btn-primary %s" style="%s">%s</button>`, classes, styleAttr, text)

	case "image":
		src := c.Content
		if src == "" {
			src = "https://via.placeholder.com/400x200"
		}
		return fmt.Sprintf(`<img src="%s" class="%s" style="%s" alt="Image" />`, src, classes, styleAttr)

	case "container":
		var childrenHTML strings.Builder
		for _, child := range c.Children {
			childrenHTML.WriteString(child.ToHTML())
		}
		return fmt.Sprintf(`<div class="p-4 %s" style="%s">%s</div>`, classes, styleAttr, childrenHTML.String())

	case "hero":
		text := c.Content
		if text == "" {
			text = "Hero Section"
		}
		return fmt.Sprintf(`<section class="hero bg-base-200 min-h-[300px] %s" style="%s">
    <div class="hero-content text-center">
        <div class="max-w-md">
            <h1 class="text-4xl font-bold mb-4">%s</h1>
            <p class="mb-8 text-base-content/70">A stunning hero section for your website.</p>
            <button class="btn btn-primary">Get Started</button>
        </div>
    </div>
</section>`, classes, styleAttr, text)

	case "navbar":
		return fmt.Sprintf(`<nav class="navbar bg-base-100 shadow-md %s" style="%s">
    <div class="flex-1">
        <a class="btn btn-ghost text-xl">Brand</a>
    </div>
    <div class="flex-none">
        <ul class="menu menu-horizontal px-1">
            <li><a>Home</a></li>
            <li><a>About</a></li>
            <li><a>Contact</a></li>
        </ul>
    </div>
</nav>`, classes, styleAttr)

	case "card":
		text := c.Content
		if text == "" {
			text = "Card Title"
		}
		return fmt.Sprintf(`<div class="card bg-base-100 shadow-md %s" style="%s">
    <div class="card-body">
        <h2 class="card-title">%s</h2>
        <p>Card content goes here.</p>
    </div>
</div>`, classes, styleAttr, text)

	case "form":
		return fmt.Sprintf(`<div class="card bg-base-100 shadow-md %s" style="%s">
    <div class="card-body">
        <form class="space-y-4">
            <div class="form-control">
                <label class="label"><span class="label-text">Name</span></label>
                <input type="text" class="input input-bordered" placeholder="Enter name" />
            </div>
            <div class="form-control">
                <label class="label"><span class="label-text">Email</span></label>
                <input type="email" class="input input-bordered" placeholder="Enter email" />
            </div>
            <button type="submit" class="btn btn-primary">Submit</button>
        </form>
    </div>
</div>`, classes, styleAttr)

	case "footer":
		return fmt.Sprintf(`<footer class="footer footer-center bg-neutral text-neutral-content p-8 %s" style="%s">
    <p>&copy; 2026 Your Company. All rights reserved.</p>
</footer>`, classes, styleAttr)

	default:
		return fmt.Sprintf(`<div class="%s" style="%s">%s</div>`, classes, styleAttr, c.Content)
	}
}

func buildStyleAttr(styles map[string]string) string {
	if styles == nil {
		return ""
	}
	var parts []string
	for k, v := range styles {
		if k == "class" {
			continue // class is handled separately
		}
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(parts, "; ")
}
