import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";
import { jwtVerify } from "jose";

export async function middleware(request: NextRequest) {
  const session = request.cookies.get("session")?.value;
  const path = request.nextUrl.pathname;

  // Define paths
  const isAuthPage = path.startsWith("/login") || path.startsWith("/register");

  // 1. If NO session
  if (!session) {
    // If trying to access protected routes, redirect to login
    if (!isAuthPage && path !== "/") {
      return NextResponse.redirect(new URL("/login", request.url));
    }
    // Otherwise, let them proceed (to login/register/landing)
    return NextResponse.next();
  }

  // 2. If Session EXISTS, verify it
  try {
    const secret = new TextEncoder().encode(
      process.env.JWT_SECRET || "test-secret-key"
    );

    await jwtVerify(session, secret);

    // ✅ Verification Successful

    // If user is on Login/Register, redirect to Dashboard
    if (isAuthPage) {
      return NextResponse.redirect(new URL("/dashboard", request.url));
    }

    // Otherwise, let them pass
    return NextResponse.next();
  } catch (error) {
    console.error("🔒 Token Invalid/Expired:", error);

    // ❌ Verification Failed (The fix for the loop)

    // Create a response depending on where they are
    let response;

    if (isAuthPage) {
      // If already on login, JUST DELETE COOKIE and let them stay
      response = NextResponse.next();
    } else {
      // If on dashboard, redirect to login
      response = NextResponse.redirect(new URL("/login", request.url));
    }

    // Nuke the bad cookie so the loop dies
    response.cookies.delete("session");

    return response;
  }
}

export const config = {
  matcher: ["/((?!api|_next/static|_next/image|favicon.ico).*)"],
};
