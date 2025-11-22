import { Skeleton } from "@/components/ui/skeleton";
import styled from "styled-components";

const Card = styled.div`
  background: #21262d;
  border: 1px solid #30363d;
  border-radius: 6px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  gap: 10px;
`;

export default function MetricsLoading() {
  return (
    <>
      {[1, 2, 3, 4].map((i) => (
        <Card key={i}>
          <Skeleton width="60%" height="12px" />
          <Skeleton width="80%" height="24px" />
        </Card>
      ))}
    </>
  );
}
