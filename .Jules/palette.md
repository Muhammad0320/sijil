## 2024-05-24 - [Accessibility] Interactive List Items
**Learning:** List items in `ProjectList` were interactive (clickable) but lacked keyboard accessibility (`tabIndex`, `role`, and key handlers). This is a common pattern in custom UI components where `div`s are used as buttons.
**Action:** Always check interactive lists for keyboard support. If using `div`s, ensure `role="button"`, `tabIndex={0}`, and `onKeyDown` (Enter/Space) are added, along with visual focus states (`:focus-visible`).
