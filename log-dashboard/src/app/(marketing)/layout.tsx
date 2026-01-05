"use client";

import styled from "styled-components";
import Link from "next/link";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";
import Footer from "@/components/layout/Footer";

import Image from "next/image";
import logo from "../../../public/logo.png";

const Nav = styled.nav`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px 48px;
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 50;
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(5, 5, 5, 0.6);
  height: 80px; /* Explicit height */
`;

const MainContent = styled.main`
  padding-top: 80px; /* MATCH THE HEADER HEIGHT */
  background-color: #050505;
  min-height: 100vh;
`;

const NavLinks = styled.div`
  display: flex;
  gap: 32px;

  a {
    font-size: 14px;
    font-weight: 500;
    color: #8b949e;
    transition: color 0.2s;
    &:hover {
      color: #fff;
    }
  }
`;

const AuthButtons = styled.div`
  display: flex;
  gap: 16px;
  align-items: center;
`;

const LogoContainer = styled(Link)`
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
`;

const LogoIcon = styled.div`
  width: 32px;
  height: 32px;
  position: relative;
`;

const BrandName = styled.span`
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.5px;
`;

export default function MarketingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div style={{ backgroundColor: "#050505", minHeight: "100vh" }}>
      <Nav>
        <LogoContainer href="/">
          <LogoIcon>
            <Image src={logo} alt="sijil logo" />
          </LogoIcon>
          <BrandName>Sijil</BrandName>
        </LogoContainer>

        {/* Hide links on mobile for now to prevent clutter */}
        <NavLinks style={{ display: "none" }}>
          {/* You can implement a proper drawer here later if needed */}
        </NavLinks>

        <AuthButtons>
          <Link
            href="/login"
            style={{ fontSize: "14px", fontWeight: 600, color: "#fff" }}
          >
            Sign In
          </Link>
          <div style={{ width: "120px" }}>
            <Link href="/register">
              <BorderBeamButton variant="primary">Start</BorderBeamButton>
            </Link>
          </div>
        </AuthButtons>
      </Nav>
      <MainContent>{children}</MainContent>
      <Footer />
    </div>
  );
}
