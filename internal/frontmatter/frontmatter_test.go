package frontmatter_test

import (
	"errors"
	"testing"

	"github.com/idoceb00/lorren/internal/frontmatter"
)

func TestSplitFrontmatter(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		wantFrontmatter string
		wantBody        string
		wantErr         error
	}{
		{
			name:            "frontmatter and body",
			input:           "---\nplan: Fuerza\n---\n\n# Título\n\nTexto.\n",
			wantFrontmatter: "plan: Fuerza",
			wantBody:        "# Título\n\nTexto.",
		},
		{
			name:            "empty body",
			input:           "---\nplan: Fuerza\n---\n",
			wantFrontmatter: "plan: Fuerza",
			wantBody:        "",
		},
		{
			name:            "closing delimiter on last line without newline",
			input:           "---\nplan: Fuerza\n---",
			wantFrontmatter: "plan: Fuerza",
			wantBody:        "",
		},
		{
			name:            "horizontal rule in body is not the closing delimiter",
			input:           "---\nplan: Fuerza\n---\n\n# A\n\n---\n\n# B\n",
			wantFrontmatter: "plan: Fuerza",
			wantBody:        "# A\n\n---\n\n# B",
		},
		{
			name:            "leading blank lines before frontmatter",
			input:           "\n\n---\nplan: Fuerza\n---\nBody\n",
			wantFrontmatter: "plan: Fuerza",
			wantBody:        "Body",
		},
		{
			name:            "windows line endings",
			input:           "---\r\nplan: Fuerza\r\n---\r\nBody\r\n",
			wantFrontmatter: "plan: Fuerza",
			wantBody:        "Body",
		},
		{
			name:    "no frontmatter",
			input:   "# Just a note\n",
			wantErr: frontmatter.ErrNotFound,
		},
		{
			name:    "empty file",
			input:   "",
			wantErr: frontmatter.ErrNotFound,
		},
		{
			name:    "unterminated frontmatter",
			input:   "---\nplan: Fuerza\n",
			wantErr: frontmatter.ErrUnterminated,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, body, err := frontmatter.Split([]byte(tt.input))

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("SplitFrontmatter() error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("SplitFrontmatter() unexpected error: %v", err)
			}
			if got := string(fm); got != tt.wantFrontmatter {
				t.Errorf("frontmatter = %q, want %q", got, tt.wantFrontmatter)
			}
			if got := string(body); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}

func TestJoinRoundTrip(t *testing.T) {
	fm := []byte("plan: Fuerza")
	body := []byte("# Título\n\nTexto.")

	gotFm, gotBody, err := frontmatter.Split(frontmatter.Join(fm, body))
	if err != nil {
		t.Fatalf("SplitFrontmatter() unexpected error: %v", err)
	}
	if string(gotFm) != string(fm) {
		t.Errorf("frontmatter = %q, want %q", gotFm, fm)
	}
	if string(gotBody) != string(body) {
		t.Errorf("body = %q, want %q", gotBody, body)
	}
}
