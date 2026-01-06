"use client";

import styled from "styled-components";
import Link from "next/link";
import { Github, Twitter, Linkedin, Activity } from "lucide-react";

const FooterContainer = styled.footer`
  background-color: #050505;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  padding: 80px 24px 40px;
  margin-top: auto; /* Push to bottom */
`;

const FooterGrid = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 48px;
  margin-bottom: 64px;

  @media (max-width: 768px) {
    grid-template-columns: 1fr 1fr;
    gap: 32px;
  }
  @media (max-width: 480px) {
    grid-template-columns: 1fr;
  }
`;

const BrandColumn = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const Logo = styled.div`
  font-size: 20px;
  font-weight: 800;
  color: #fff;
  display: flex;
  align-items: center;
  gap: 8px;

  &::before {
    content: "";
    width: 10px;
    height: 10px;
    background: #238636;
    border-radius: 50%;
    box-shadow: 0 0 10px rgba(35, 134, 54, 0.5);
  }
`;

const Tagline = styled.p`
  color: #8b949e;
  font-size: 14px;
  line-height: 1.6;
  max-width: 300px;
`;

const LinkColumn = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

const ColumnTitle = styled.h4`
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
`;

const FooterLink = styled(Link)`
  color: #8b949e;
  font-size: 14px;
  transition: color 0.2s;

  &:hover {
    color: #58a6ff;
  }
`;

const BottomBar = styled.div`
  max-width: 1200px;
  margin: 0 auto;
  padding-top: 32px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  justify-content: space-between;
  align-items: center;
  flex-wrap: wrap;
  gap: 16px;

  @media (max-width: 768px) {
    flex-direction: column;
    text-align: center;
  }
`;

const Copyright = styled.span`
  color: #484f58;
  font-size: 13px;
`;

const Socials = styled.div`
  display: flex;
  gap: 24px;

  a {
    color: #8b949e;
    transition: color 0.2s;
    &:hover {
      color: #fff;
    }
  }
`;

const StatusBadge = styled.div`
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: #2ecc71;
  background: rgba(46, 204, 113, 0.1);
  padding: 4px 8px;
  border-radius: 12px;
  border: 1px solid rgba(46, 204, 113, 0.2);
  width: fit-content;
`;

export default function Footer() {
  return (
    <FooterContainer>
      <FooterGrid>
        <BrandColumn>
          <Logo>Sijil</Logo>
          <Tagline>
            The observability platform designed for the AI era. Ingest, Query,
            and visualize at hyperspeed.
          </Tagline>
          <div style={{ marginTop: "auto" }}>
            <StatusBadge>
              <Activity size={12} />
              <span>All Systems Operational</span>
            </StatusBadge>
          </div>
        </BrandColumn>

        <LinkColumn>
          <ColumnTitle>Product</ColumnTitle>
          <FooterLink href="#features">Architecture</FooterLink>
          <FooterLink href="#performance">Benchmarks</FooterLink>
          <FooterLink href="#pricing">Pricing</FooterLink>
          <FooterLink href="/changelog">Changelog</FooterLink>
        </LinkColumn>

        <LinkColumn>
          <ColumnTitle>Resources</ColumnTitle>
          <FooterLink href="/docs">Documentation</FooterLink>
          <FooterLink href="/api">API Reference</FooterLink>
          <FooterLink href="/blog">Engineering Blog</FooterLink>
          <FooterLink href="/community">Community</FooterLink>
        </LinkColumn>

        <LinkColumn>
          <ColumnTitle>Legal</ColumnTitle>
          <FooterLink href="/privacy">Privacy Policy</FooterLink>
          <FooterLink href="/terms">Terms of Service</FooterLink>
          <FooterLink href="/security">Security</FooterLink>
        </LinkColumn>
      </FooterGrid>

      <BottomBar>
        <Copyright>
          &copy; {new Date().getFullYear()} Sijil Inc. All rights reserved.
        </Copyright>
        <Socials>
          <Link
            href="https://github.com/muhammad0320"
            target="_blank"
            aria-label="GitHub"
          >
            <Github size={18} />
          </Link>
          <Link href="https://twitter.com" target="_blank" aria-label="Twitter">
            <Twitter size={18} />
          </Link>
          <Link
            href="https://linkedin.com"
            target="_blank"
            aria-label="LinkedIn"
          >
            <Linkedin size={18} />
          </Link>
        </Socials>
      </BottomBar>
    </FooterContainer>
  );
}
