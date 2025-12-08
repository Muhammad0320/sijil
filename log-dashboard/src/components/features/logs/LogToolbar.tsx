"use client";

import styled from "styled-components";
import { Search, X } from "lucide-react";

const Container = styled.div`
  height: 50px;
  border-bottom: 1px solid #30363d;
  display: flex;
  align-items: center;
  padding: 0 16px;
  background: #0d1117;
  gap: 12px;
`;

const SearchContainer = styled.div`
  position: relative;
  width: 300px;
`;

const SearchInput = styled.input`
  width: 100%;
  background: #161b22;
  border: 1px solid #30363d;
  border-radius: 6px;
  padding: 6px 10px 6px 32px;
  color: #c9d1d9;
  font-size: 13px;

  &:focus {
    outline: none;
    border-color: #58a6ff;
  }
`;

const IconWrapper = styled.div`
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  color: #8b949e;
  display: flex;
`;

interface LogToolbarProps {
  searchQuery: string;
  setSearchQuery: (q: string) => void;
  onRefresh: () => void; // For manual refresh if needed
}

export function LogToolbar({
  searchQuery,
  setSearchQuery,
  onRefresh,
}: LogToolbarProps) {
  return (
    <Container>
      <SearchContainer>
        <IconWrapper>
          <Search size={14} />
        </IconWrapper>
        <SearchInput
          placeholder="Filter logs..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </SearchContainer>

      {/* Add Time Range Picker here later */}
    </Container>
  );
}
