import React from 'react'
import styled from "@emotion/styled";
import {useDocumentTitle} from "@/utils";

export const ServiceAgreement = () =>{
    useDocumentTitle("念熹智语产品服务协议");
   return <Container>
        <Section  >
            <Content>
                <h2 style={{textAlign: 'center'}}>念熹智语产品服务协议</h2> <p>欢迎您使用念熹智语（NX)产品。</p>
                    <p>请您仔细阅读以下条款，如果您对本协议的任何条款表示异议，您可以选择不使用念熹智语。鉴于《念熹科技用户协议》已经包含本《念熹智语产品服务协议》（“《念熹智语服务协议》”）及其更新版本（念熹科技会及时提示您更新的情况），您同意和接受《念熹科技用户协议》即表明您已同意本《念熹智语产品服务协议》（含更新版本）的条款并受其约束。</p>
                    <h4>1.免责声明</h4> <p>
                        1.1.念熹智语自动生成中文写作的修改建议和结果（“修改成果”），修改成果没有经过任何人工整理与编辑。
                    </p> <p>
                        1.2.念熹智语仅为用户进行中文写作所使用的工具。用户使用念熹智语的任何服务不应当视为用户与念熹科技之间就修改成果成立任何委托、承揽、加工等民事法律关系。
                    </p> <p>
                        1.3.念熹智语对修改成果的正确性和合法性不做任何形式的保证，亦不承担任何法律责任。用户对其使用修改成果以及对修改成果做出的修改造成的任何后果独自承担全部责任。
                    </p> <p>
                        1.4.念熹智语根据用户键入的内容及上传的文件自动生成修改成果，不代表念熹科技赞成键入内容、上传文件及修改成果的内容或立场。
                    </p> <p>
                        1.5.念熹智语仅根据用户键入内容及上传文件提供机器修改成果。就用户使用念熹智语过程中键入的内容、上传的文件、修改成果以及用户对于修改成果做出的修改，用户已获得任何第三方的必要和合法授权。念熹科技不对用户键入的内容、上传的文件、修改成果以及用户对于修改成果做出的修改承担任何法律责任。
                    </p> <h4>2.保密承诺</h4> <p>
                        2.1.念熹科技承诺对用户使用念熹智语过程中键入的内容、上传的文件以及修改成果进行保密，不向任何第三方提供、披露其中的信息，也不会将其中信息用于非本协议之目的。
                    </p> <p>
                        2.2.上述保密信息不包括：1）非由念熹科技的违约行为而导致公开的信息；2）接收到信息之前，既已属于念熹科技的信息；或2）念熹科技通过合法渠道从第三方获得的信息。
                    </p> <p>
                        2.3.您了解并同意，念熹科技可以将您对于修改成果的任何修改及其相关数据用于提升念熹智语的系统升级和算法修正。念熹科技不会向任何第三方提供、披露您的修改及其相关数据。
                    </p> <p>
                        2.4.如果您希望删除键入的内容、上传的文件、修改成果以及您对于修改成果做出的修改，您可以通过页面提供的选项选择删除上述内容。在您选择删除后，我们不继续保留相关信息。
                    </p> <h4>3.修订和通知</h4> <p>
                        由于服务内容、法律和监管政策要求等变化，我们可能会适时对本念熹智语服务协议进行修订。当本念熹智语服务协议发生变更时，我们会在我们的平台发布或向您提供的电子邮箱地址发送更新后的版本。为了您能及时接收到通知，建议您在电子邮箱地址变动时及时更新账号信息或通知我们。如您继续使用服务，表示同意接受修订后的本念熹智语服务协议的内容。
                    </p>
            </Content>
        </Section>
    </Container>
}


const Container = styled.div`
  width: 60vw;
  line-height: 25px;
  overflow-y: auto;
  margin: 0 auto;
  background-color: #fff;
  
`

const Section = styled.section`
  display: -webkit-box;
  display: -ms-flexbox;
  display: flex;
  -webkit-box-orient: horizontal;
  -webkit-box-direction: normal;
  -ms-flex-direction: row;
  flex-direction: row;
  -webkit-box-flex: 1;
  -ms-flex: 1;
  flex: 1;
  -ms-flex-preferred-size: auto;
  flex-basis: auto;
  -webkit-box-sizing: border-box;
  box-sizing: border-box;
  min-width: 0;
  overflow: auto;
`

const Content = styled.div`

`
const UL = styled.ul`
  list-style-type: none;
  margin:0px;
  padding: 0px;
`
