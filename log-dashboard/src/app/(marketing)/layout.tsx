"use client";

import styled from "styled-components";
import Link from "next/link";
import { BorderBeamButton } from "@/components/ui/borderBeamButton";
import Footer from "@/components/layout/Footer";

const Nav = styled.nav`
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 48px;
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 50;
  backdrop-filter: blur(12px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
  background: rgba(5, 5, 5, 0.6);
  height: 80px;
  transition: all 0.2s ease-in-out;

  @media (max-width: 768px) {
    padding: 0 24px;
  }
`;

const MainContent = styled.main`
  padding-top: 80px;
  background-color: #050505;
  min-height: 100vh;
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
  background: #58a6ff;
  border-radius: 6px;
  position: relative;
  overflow: hidden;

  /* Abstract shutter/speed lines */
  &::after {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    width: 200%;
    height: 4px;
    background: #000;
    transform: translate(-50%, -50%) rotate(45deg);
  }
`;

const BrandName = styled.span`
  font-size: 24px;
  font-weight: 800;
  letter-spacing: -1px;
  color: #fff;
`;

const NavLinks = styled.div`
  display: flex;
  gap: 32px;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);

  a {
    font-size: 14px;
    font-weight: 500;
    color: #8b949e;
    transition: color 0.2s;
    &:hover {
      color: #fff;
    }
  }

  @media (max-width: 768px) {
    display: none;
  }
`;

const AuthButtons = styled.div`
  display: flex;
  gap: 16px;
  align-items: center;
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
          <LogoIcon />
          <BrandName>Sijil</BrandName>
        </LogoContainer>

        <NavLinks>
          <Link href="#features">Features</Link>
          <Link href="#performance">Benchmarks</Link>
          <Link href="#pricing">Pricing</Link>
          <Link href="/docs">Docs</Link>
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
