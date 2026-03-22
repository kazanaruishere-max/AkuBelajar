# 🎨 UI Component Specification — AkuBelajar

> Design tokens, component taxonomy, accessibility, responsiveness, dan dark mode.

---

## 1. Design Tokens

### Colors

```css
:root {
  /* Primary — Blue (Brand) */
  --color-primary-50:  #EFF6FF;
  --color-primary-100: #DBEAFE;
  --color-primary-200: #BFDBFE;
  --color-primary-500: #3B82F6;
  --color-primary-600: #2563EB;
  --color-primary-700: #1D4ED8;
  --color-primary-900: #1E3A5F;

  /* Success — Green */
  --color-success-50:  #F0FDF4;
  --color-success-500: #22C55E;
  --color-success-700: #15803D;

  /* Warning — Amber */
  --color-warning-50:  #FFFBEB;
  --color-warning-500: #F59E0B;
  --color-warning-700: #B45309;

  /* Danger — Red */
  --color-danger-50:  #FEF2F2;
  --color-danger-500: #EF4444;
  --color-danger-700: #B91C1C;

  /* Neutral — Slate */
  --color-neutral-50:  #F8FAFC;
  --color-neutral-100: #F1F5F9;
  --color-neutral-200: #E2E8F0;
  --color-neutral-300: #CBD5E1;
  --color-neutral-500: #64748B;
  --color-neutral-700: #334155;
  --color-neutral-800: #1E293B;
  --color-neutral-900: #0F172A;

  /* Background */
  --bg-primary: #FFFFFF;
  --bg-secondary: #F8FAFC;
  --bg-tertiary: #F1F5F9;
}
```

### Typography

```css
:root {
  --font-sans: 'Inter', -apple-system, sans-serif;
  --font-mono: 'JetBrains Mono', monospace;

  /* Font sizes (rem) */
  --text-xs:   0.75rem;    /* 12px */
  --text-sm:   0.875rem;   /* 14px */
  --text-base: 1rem;       /* 16px */
  --text-lg:   1.125rem;   /* 18px */
  --text-xl:   1.25rem;    /* 20px */
  --text-2xl:  1.5rem;     /* 24px */
  --text-3xl:  1.875rem;   /* 30px */
  --text-4xl:  2.25rem;    /* 36px */

  /* Font weights */
  --font-normal:   400;
  --font-medium:   500;
  --font-semibold: 600;
  --font-bold:     700;
}
```

### Spacing & Sizing

```css
:root {
  --space-1: 4px;
  --space-2: 8px;
  --space-3: 12px;
  --space-4: 16px;
  --space-5: 20px;
  --space-6: 24px;
  --space-8: 32px;
  --space-10: 40px;

  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --radius-xl: 16px;
  --radius-full: 9999px;

  --shadow-sm: 0 1px 2px rgba(0,0,0,0.05);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.07);
  --shadow-lg: 0 10px 15px rgba(0,0,0,0.1);
}
```

---

## 2. Component Taxonomy

### Atoms (Primitif)

| Component | Varian | Usage |
|:---|:---|:---|
| `Button` | primary, secondary, outline, ghost, danger | CTA, actions |
| `Input` | text, email, password, number, textarea | Form input |
| `Select` | single, multi, searchable | Dropdown |
| `Checkbox` | default, indeterminate | Multiple choice |
| `RadioGroup` | default | Single choice |
| `Badge` | success, warning, danger, info, neutral | Status label |
| `Avatar` | sm (32), md (40), lg (48) | User photo |
| `Spinner` | sm, md, lg | Loading indicator |
| `Skeleton` | text, card, table-row | Loading placeholder |

### Molecules (Gabungan)

| Component | Komposisi | Usage |
|:---|:---|:---|
| `FormField` | Label + Input + ErrorText | Form rows |
| `SearchInput` | Input + Icon + ClearButton | Filter/search |
| `DatePicker` | Input + Calendar Popover | Date selection |
| `TimePicker` | Input + TimeDial | Time selection |
| `FileUpload` | DropZone + ProgressBar + FileList | Upload files |
| `CameraCapture` | Video + ShutterBtn + Preview + Toggle | 📷 Foto tugas (NO MIRROR, rear default) |
| `QRScanner` | Video + Overlay + html5-qrcode | 📱 Scan QR absensi (NO MIRROR) |
| `Toast` | Icon + Message + CloseButton | Notifications |
| `Modal` | Overlay + Card + Actions | Dialogs |
| `ConfirmDialog` | Modal + Warning text + Cancel/Confirm | Destructive confirmations |
| `DataTable` | Table + Pagination + Sort + Filter | List data |
| `EmptyState` | Icon + Title + Description + Action | No data |

### Organisms (Page-Level)

| Component | Komposisi | Usage |
|:---|:---|:---|
| `Sidebar` | Logo + NavLinks + UserMenu | Dashboard nav |
| `TopBar` | Breadcrumb + NotifBell + Avatar | Header |
| `PageHeader` | Title + Description + Actions | Page title area |
| `StatsCard` | Icon + Number + Label + Trend | KPI display |
| `QuizPlayer` | Timer + QuestionCard + Navigation | CBT interface |

---

## 3. Responsive Breakpoints

| Breakpoint | Width | Device | Layout |
|:---|:---|:---|:---|
| `sm` | ≥ 375px | Mobile | 1 column, bottom nav |
| `md` | ≥ 768px | Tablet | 2 columns, side nav collapsed |
| `lg` | ≥ 1024px | Laptop | 3 columns, side nav open |
| `xl` | ≥ 1280px | Desktop | Full layout |
| `2xl` | ≥ 1536px | Large desktop | Max-width container |

### Layout Rules

- Dashboard: Sidebar (240px fixed) + Content (fluid)
- Mobile: Sidebar → Bottom tab bar
- Tables: Horizontal scroll di mobile (< 768px)
- Forms: Single column di mobile, 2 kolom di desktop

---

## 4. Accessibility (a11y) Checklist

| Rule | Detail |
|:---|:---|
| Color contrast | WCAG AA minimum (4.5:1 normal text, 3:1 large text) |
| Focus visible | Semua interactive element harus ada focus ring |
| Keyboard nav | Tab order logis, Enter/Space untuk aksi |
| Screen reader | ARIA labels di semua icon-only buttons |
| Error messages | `aria-describedby` link ke error text |
| Form labels | Setiap input harus punya `<label>` |
| Alt text | Semua `<img>` harus punya alt text deskriptif |
| Motion | `prefers-reduced-motion` → disable animasi |
| Language | `<html lang="id">` |

---

## 5. Dark Mode

```css
[data-theme="dark"] {
  --bg-primary: #0F172A;
  --bg-secondary: #1E293B;
  --bg-tertiary: #334155;
  
  --color-neutral-50: #0F172A;
  --color-neutral-100: #1E293B;
  --color-neutral-700: #CBD5E1;
  --color-neutral-800: #E2E8F0;
  --color-neutral-900: #F8FAFC;
  
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.3);
  --shadow-md: 0 4px 6px rgba(0,0,0,0.4);
}
```

- Toggle: `localStorage` preference → Zustand uiStore
- Default: follow `prefers-color-scheme` OS
- Transition: `transition: background-color 0.2s, color 0.2s`

---

*Terakhir diperbarui: 22 Maret 2026*
