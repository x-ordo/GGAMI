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
		return fmt.Sprintf(`<h1 class="%s" style="%s">%s</h1>`, classes, styleAttr, text)

	case "paragraph":
		text := c.Content
		if text == "" {
			text = "Paragraph text here..."
		}
		return fmt.Sprintf(`<p class="%s" style="%s">%s</p>`, classes, styleAttr, text)

	case "button":
		text := c.Content
		if text == "" {
			text = "Click Me"
		}
		return fmt.Sprintf(`<button class="%s" style="%s">%s</button>`, classes, styleAttr, text)

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
		return fmt.Sprintf(`<div class="%s" style="%s">%s</div>`, classes, styleAttr, childrenHTML.String())

	case "hero":
		text := c.Content
		if text == "" {
			text = "Hero Section"
		}
		return fmt.Sprintf(`<section class="bg-gray-900 text-white py-20 px-10 text-center rounded-xl %s" style="%s">
    <h1 class="text-4xl font-bold mb-4">%s</h1>
    <p class="text-xl text-gray-400 mb-8">A stunning hero section for your website.</p>
    <button class="bg-yellow-500 text-black font-bold py-3 px-8 rounded-full hover:bg-yellow-400 transition">Get Started</button>
</section>`, classes, styleAttr, text)

	case "navbar":
		return fmt.Sprintf(`<nav class="bg-white shadow-md px-6 py-4 flex items-center justify-between %s" style="%s">
    <span class="text-xl font-bold">Brand</span>
    <div class="flex gap-4">
        <a href="#" class="text-gray-600 hover:text-gray-900">Home</a>
        <a href="#" class="text-gray-600 hover:text-gray-900">About</a>
        <a href="#" class="text-gray-600 hover:text-gray-900">Contact</a>
    </div>
</nav>`, classes, styleAttr)

	case "card":
		text := c.Content
		if text == "" {
			text = "Card Title"
		}
		return fmt.Sprintf(`<div class="bg-white rounded-lg shadow-md p-6 %s" style="%s">
    <h3 class="text-lg font-bold mb-2">%s</h3>
    <p class="text-gray-600">Card content goes here.</p>
</div>`, classes, styleAttr, text)

	case "form":
		return fmt.Sprintf(`<form class="bg-white p-6 rounded-lg shadow-md space-y-4 %s" style="%s">
    <div>
        <label class="block text-gray-700 font-medium mb-1">Name</label>
        <input type="text" class="w-full border rounded p-2" placeholder="Enter name" />
    </div>
    <div>
        <label class="block text-gray-700 font-medium mb-1">Email</label>
        <input type="email" class="w-full border rounded p-2" placeholder="Enter email" />
    </div>
    <button type="submit" class="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">Submit</button>
</form>`, classes, styleAttr)

	case "footer":
		return fmt.Sprintf(`<footer class="bg-gray-800 text-gray-400 py-8 px-6 text-center %s" style="%s">
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
