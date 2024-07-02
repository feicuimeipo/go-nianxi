import React from "react";
import logo from "@/assets/logo.png";
import styled from "@emotion/styled";

export const LeftNavLogo = ({ subtitle }: { subtitle: string}) => {
    return <Container>
        <Logo>
            {subtitle? <Title>{subtitle}</Title>:null}
        </Logo>
    </Container>
};

const Container = styled.div`   
`

const Logo = styled.div`
  background: url(${logo}) no-repeat left;
  background-position-x: 30px;  
  //上右下左
  //如果提供两个，第一个用于上、下，第二个用于左、右。
  padding: 1rem 0rem;
  background-size: 10rem;
  vertical-align: bottom;  
`;

const Title = styled.h3`
  padding: 8rem 0rem 0rem 3.5rem;
  margin-bottom: 2rem;
  color: rgb(84, 84, 85);
  font-size: 2rem;
`;
