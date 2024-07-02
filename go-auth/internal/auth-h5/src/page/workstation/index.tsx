import React from "react";
import styled from "@emotion/styled";
import WorkstationData from "@/page/workstation/data";
import {Outlet} from "react-router";
import {NavMenuBar} from "@/components/menu/menu";
import "./index.css";
import {Button, message} from "antd";
import {useAsync} from "@/utils/use-async";
import {useAuth} from "@/context/auth-context";
import {LeftNavLogo} from "@/components/logo";


export const WorkstationApp = () => {

    const { run, isLoading } = useAsync(undefined, { throwOnError: true });
    const {logout} = useAuth()

    const handleLogout = async (values) => {
        try {
            await run(logout())
        } catch (e: any) {
            message.error(e,1000)
        }
    }


    return  <Container>
                <Nav>
                    <LeftNavLogo  subtitle={"用户中心"} />
                    <NavMenuBar  menuList={WorkstationData.navMenuList}/>
                </Nav>
                <Header>
                    <HeaderLeft></HeaderLeft>
                    <HeaderCenter></HeaderCenter>
                    <HeaderRight>
                        <Button onClick={handleLogout}  loading={isLoading} htmlType={"button"} >退出</Button>

                    </HeaderRight>
                </Header>
                <Main>
                    <Outlet />
                </Main>
           </Container>
}

const Container = styled.div`
  display: grid;
  grid-template-columns: var(--home-nav-panel-width) auto ;
  grid-template-rows: var(--home-header-panel-height) auto ;
  grid-template-areas:
    'nav header '
    'nav main '
    'nav main ';
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  height: 100vh;
`

const Nav = styled.nav`
  grid-area: nav;
  border-right: 1px solid var(--home-border-color);  
`

const Header = styled.header`
  grid-area: header;
  border-bottom: 1px solid var(--home-border-color);
  vertical-align: middle;
  padding: 0.5rem;  
`

const Main = styled.main`
  grid-area: main;  
`

// const Logo = styled.div`
//   background: url(${logo}) no-repeat left;
//   background-position-x: 30px;
//   padding: 5rem 0rem 5rem 5rem;
//   background-size: 8rem;
//   vertical-align: bottom;
// `;


const HeaderLeft = styled.div`
  float: left;
`

const HeaderCenter = styled.div``

const HeaderRight = styled.div`
  float: right;
`


export default WorkstationApp
