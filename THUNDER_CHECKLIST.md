## Thunder API

* [X] Get
	- [X] Query
* [X] Patch
	- [X] SObject
* [ ] Post
	- [X] SObject
	- [ ] Composite
* [X] Delete
	- [X] SObject
* [x] RecordId

- [X] Query
	- [x] Improve Record Access API

## Thunder SLDS Components Roadmap

### Currently Implemented Components ✅

#### Form Components
- [x] **TextInput** (`input.go`) - Basic text input with label
- [x] **Textarea** (`textarea.go`) - Multi-line text input
- [x] **Checkbox** (`checkbox.go`) - Boolean input 
- [x] **Select** (`select.go`) - Dropdown selection
- [x] **RadioGroup** (`radiogroup.go`) - Single-choice options with form layout
- [x] **RadioButtonGroup** (`radiogroup.go`) - Single-choice options with button styling
- [x] **Datepicker** (`datepicker.go`) - Date selection
- [x] **Lookup** (`lookup.go`) - Autocomplete/search input
- [x] **Form Validation Framework** (`validated_input.go`) - Field-level validation with ValidationState

#### Navigation Components  
- [x] **Breadcrumb** (`breadcrumb.go`) - Navigation hierarchy
- [x] **Tabs** (`tabs.go`) - Tabbed navigation with content panels

#### Data Display Components
- [x] **DataTable** (`datatable.go`) - Basic tabular display
- [x] **ProgressBar** (`progressbar.go`) - Horizontal progress indicator
- [x] **Badge** (`badge.go`) - Status/label display
- [x] **Icon** (`icon.go`) - SLDS icon rendering

#### Layout Components
- [x] **Card** (`card.go`) - Container with header and body
- [x] **Grid** (`grid.go`) - SLDS grid system
- [x] **Page** (`page.go`) - Layout wrapper
- [x] **PageHeader** (`pageheader.go`) - Page-level heading

#### Feedback Components
- [x] **Modal** (`modal.go`) - Overlay dialogs
- [x] **Toast** (`toast.go`) - Notification messages
- [x] **Spinner** (`spinner.go`) - Loading indicators
- [x] **Stencil** (`stencil.go`) - Skeleton loading placeholders
- [x] **Tooltip** (`validated_input.go`) - Contextual help with info icons

#### Utility Components
- [x] **Button** (`button.go`) - Action buttons (neutral, brand, destructive)

### Critical Missing Components (High Priority) 🔴

#### Form Components & Validation
1. **File Upload** - File selection and upload
   - *Priority: MEDIUM* - Important for data import/export scenarios
   - *Effort: MEDIUM* - File handling and progress display
   - *Usage: Common* - Document attachments, CSV imports

2. **Combobox** - Enhanced dropdown with search/filtering
   - *Priority: MEDIUM* - More advanced than basic Select
   - *Effort: MEDIUM* - Builds on Lookup functionality
   - *Usage: Common* - Large option lists, picklist values

#### Navigation Components
3. **Vertical Navigation** - Sidebar/tree navigation
   - *Priority: HIGH* - Essential for complex app navigation
   - *Effort: MEDIUM* - Tree structure with expand/collapse
   - *Usage: High* - App sidebars, hierarchical menus

4. **Menu/Dropdown** - Context menus and action dropdowns
   - *Priority: HIGH* - Critical for action-heavy interfaces
   - *Effort: MEDIUM* - Positioning and click-outside handling
   - *Usage: High* - Table actions, button groups

#### Data Display Components
5. **Tree Grid** - Hierarchical data display
   - *Priority: MEDIUM* - Important for complex data relationships
   - *Effort: HIGH* - Complex tree structure with DataTable features
   - *Usage: Medium* - Folder structures, org hierarchies

6. **Accordion** - Collapsible content sections
   - *Priority: MEDIUM* - Good for organizing content
   - *Effort: LOW* - Expand/collapse state management
   - *Usage: Medium* - FAQ sections, grouped settings

#### Layout Components
7. **Tiles** - Grid-based content layout
   - *Priority: LOW* - Nice-to-have for dashboard layouts
   - *Effort: LOW* - Extension of Grid system
   - *Usage: Medium* - Dashboard widgets, gallery views

#### Feedback Components
8. **Prompt/Confirmation Dialog** - User decision dialogs
   - *Priority: HIGH* - Critical for destructive actions
   - *Effort: LOW* - Extension of Modal component
   - *Usage: High* - Delete confirmations, unsaved changes

### Nice-to-Have Components (Medium Priority) 🟡

#### Advanced Form Components
9. **DateTime Picker** - Combined date and time selection
   - *Effort: MEDIUM* - Combines existing Datepicker with time
   - *Usage: Medium* - Scheduling, event creation

10. **Color Picker** - Color selection input
    - *Effort: MEDIUM* - Color wheel/palette implementation
    - *Usage: Low* - Theming, customization features

11. **Rich Text Editor** - WYSIWYG text editing
    - *Effort: HIGH* - Complex text formatting capabilities
    - *Usage: Medium* - Email composition, content editing

#### Data Display Components  
12. **Progress Ring** - Circular progress indicator
    - *Effort: LOW* - SVG-based circular progress
    - *Usage: Low* - Alternative to ProgressBar

13. **Carousel** - Image/content slideshow
    - *Effort: MEDIUM* - Touch/swipe support, navigation
    - *Usage: Low* - Product galleries, onboarding

#### Specialized Components
14. **Avatar** - User profile images
    - *Effort: LOW* - Image with fallback initials
    - *Usage: Medium* - User lists, profiles

15. **Button Group** - Grouped action buttons
    - *Effort: LOW* - Styling for button collections
    - *Usage: Medium* - Toolbar actions, toggle groups

### Implementation Strategy & Elm Architecture Fit 🟢

#### Phase 1: Critical Form & Validation (Sprint 1-2) ✅ COMPLETED
- [x] **Textarea** - Extends existing input patterns
- [x] **Form Validation Framework** - Msg-based validation state with ValidationState
- [x] **Tooltip** - Essential UX improvement with help icons

#### Phase 2: Navigation & Menus (Sprint 3-4)  
- **Vertical Navigation** - Tree-like Elm Architecture model
- **Menu/Dropdown** - Click handling with outside detection
- **Prompt/Confirmation** - Modal variant for decisions

#### Phase 3: Advanced Data Display (Sprint 5-6)
- **Accordion** - Expand/collapse state management
- **Tree Grid** - Hierarchical data with CRUD operations
- **File Upload** - Progress tracking via Elm commands

#### Phase 4: Polish & Enhancement (Sprint 7+)
- **Combobox** - Enhanced Select with filtering
- **Avatar**, **Button Group**, **Tiles** - UI polish
- **Rich Text Editor** - Complex state management

### Elm Architecture Considerations

#### State Management Patterns
- [x] **Form Validation**: Field-level validation state in model (ValidationState struct)
- **Navigation**: Route/selection state with Msg routing  
- **Tree Components**: Recursive data structures with expand state
- **File Upload**: Progress tracking via Cmd/Msg patterns
- [x] **Tooltips**: Simple title attribute approach (no complex positioning state needed)

#### Message Patterns
- [x] `ValidationMsg { Field, Error }` - Field validation results (implemented in examples)
- `NavigationMsg { Route, Action }` - Navigation state changes  
- `TreeToggleMsg { NodeId, Expanded }` - Tree node expansion
- `FileUploadMsg { Progress, Status }` - Upload progress updates

#### Command Patterns
- [x] Async validation with debounced inputs (demonstrated in validation example)
- File upload progress tracking
- Navigation route changes
- Data fetching for combobox options

## Thunder CLI Features

### Deploy

- [x] Add CustomTab to package.xml
- [x] Display deploy errors
- [x] Sanitize App Names
