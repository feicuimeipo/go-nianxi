import React from "react";
import styled from "@emotion/styled";

const Agreement = () =>{
    const serviceAgreement = () =>{
        window.open(`/serviceAgreement`, '_blank');
    }

    const customerAgreement = () =>{
        window.open(`/customerAgreement`, '_blank');
    }
    return <Container><a href={"#"} type={"link"} onClick={customerAgreement}>用户协议</a>与<a href={"#"}  type={"link"} onClick={serviceAgreement}>产品服务协议</a></Container>
}

export default Agreement

const Container = styled.div`
    width: 100%;
    padding-top: 10px;
`
