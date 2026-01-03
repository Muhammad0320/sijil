from playwright.sync_api import sync_playwright, expect
import jwt
import datetime

def test_dashboard_accessibility(page):
    # 1. Forge a JWT
    secret = "test-secret-key"
    token = jwt.encode({
        "userId": 1,
        "exp": datetime.datetime.now(datetime.timezone.utc) + datetime.timedelta(hours=1)
    }, secret, algorithm="HS256")

    # 2. Set the cookie
    context = page.context
    context.add_cookies([{
        "name": "session",
        "value": token,
        "domain": "localhost",
        "path": "/",
        "httpOnly": True,
        "secure": False,
        "sameSite": "Lax"
    }])

    # 3. Go to Dashboard
    page.goto("http://localhost:3000/dashboard")

    # 4. Wait for Projects header to verify we are in (or at least loaded the client component)
    try:
        page.wait_for_selector("text=Projects", timeout=5000)
    except:
        print("Could not find 'Projects' text. Maybe stuck on login or error?")
        # Take a screenshot anyway to see what happened
        page.screenshot(path="verification/verification_failed.png")
        return

    # 5. Verify ARIA label on Add Button
    # The AddButton is likely an icon-only button inside the header.
    # We can search by role and name (aria-label).
    add_button = page.get_by_role("button", name="Create new project")

    # Assert it is visible
    if add_button.is_visible():
        print("SUCCESS: Found button with aria-label='Create new project'")
        # Highlight it for the screenshot
        add_button.highlight()
    else:
        print("FAILURE: Button with aria-label='Create new project' NOT found.")

    # 6. Screenshot
    page.screenshot(path="verification/verification.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_dashboard_accessibility(page)
        finally:
            browser.close()
