import React, { Component } from 'react'
import { getPicture, reqCheck} from '@/api/captcha/index'
import '@/assets/verify.css';
import defaultImg from '@/assets/default.jpg'
import {aesEncrypt} from "@/api/ase";

class VerifyPointFixed extends Component {
  constructor(props) {
    super(props);
    this.state = {
      secretKey: '', //后端返回的ase加密秘钥
      checkNum: 3, //默认需要点击的字数
      fontPos: [], //选中的坐标信息
      checkPosArr: [], //用户点击的坐标
      num: 1, //点击的记数
      pointBackImgBase: '', //后端获取到的背景图片
      poinTextList: [], //后端返回的点击字体顺序
      backToken: '', //后端返回的token值
      captchaType: 'clickWord',
      setSize: {
        imgHeight: 0,
        imgWidth: 0,
        barHeight: 0,
        barWidth: 0,
      },
      tempPoints: [],
      text: '',
      barAreaColor: 'rgb(0,0,0)',
      barAreaBorderColor: 'rgb(221, 221, 221)',
      showRefresh: true,
      bindingClick: true
    };
  }
  componentDidMount() {
    this.uuid()
    this.getData()
  }
  // 初始话 uuid 
  uuid() {
    const s = [];
    const hexDigits = "0123456789abcdef";
    for (let i = 0; i < 36; i++) {
        s[i] = hexDigits.substr(Math.floor(Math.random() * 0x10), 1);
    }
    s[14] = "4";
    s[19] = hexDigits.substr((s[19] & 0x3) | 0x8, 1); // bits 6-7 of the clock_seq_hi_and_reserved to 01
    s[8] = s[13] = s[18] = s[23] = "-";
    const slider = 'slider-' + s.join("");
    const point = 'point-'+ s.join("");
    // 判断下是否存在 slider
    if(!localStorage.getItem('slider')) {
      localStorage.setItem('slider', slider)
    }
    if(!localStorage.getItem('point')) {
      localStorage.setItem("point",point);
    }
  }

  // 初始化数据
  getData() {
    const params = {captchaType: this.state.captchaType,
                  clientUid: localStorage.getItem('point'),
                  ts: Date.now()}
    getPicture(params).then(res => {
      if(res.code === "0") {
        this.setState({
          pointBackImgBase: res.data.originalImageBase64,
          backToken: res.data.token,
          secretKey: res.data.secretKey,
          text: '请依次点击【' + res.data.wordList.join(",") + '】'
        })
      }else{
        this.setState({
          // 请求次数超限
          pointBackImgBase: null,
          text: res.msg,
          barAreaColor: '#d9534f',
          barAreaBorderColor: '#d9534f'
        })
      }
    })
  }

  // 刷新
  refresh = () => {
    this.getData()
    this.setState({
      num: 1,
      tempPoints: [],
      bindingClick: true,
      barAreaColor: 'rgb(0,0,0)',
      barAreaBorderColor: 'rgb(221, 221, 221)',
    })
    this.props.verifyPointFixedChild(null)
  }
  
  canvasClick = (e) => {
    if(this.state.bindingClick) {
      this.state.tempPoints.push(this.getMousePos(e))
      this.setState({
          tempPoints: this.state.tempPoints
      })
      if(this.state.num === this.state.checkNum) {
        this.setState({
           bindingClick: false
        })
        let data:CaptchaCheckRequest  = {
            captchaType:this.state.captchaType,
            pointJson:this.state.secretKey? aesEncrypt(JSON.stringify(this.state.tempPoints),this.state.secretKey):JSON.stringify(this.state.tempPoints),
            token:this.state.backToken,
            clientUid: localStorage.getItem('point'),
            ts: Date.now()
        }
        reqCheck(data).then(res => {
          if(res.code === "0") {
              this.setState({
                  text: '验证成功',
                  barAreaColor: '#4cae4c',
                  barAreaBorderColor: '#5cb85c'
              })
              setTimeout(() => {
                this.closeBox();
              }, 1500)
            this.props.verifyPointFixedChild(data)
          } else {
              this.setState({
                text: res.msg,
                barAreaColor: '#d9534f',
                barAreaBorderColor: '#d9534f'
              })
              setTimeout(() => {
                  this.refresh();
              }, 1000);

          }
        }).catch(res=>{
          this.setState({
            text: res.msg,
            barAreaColor: '#d9534f',
            barAreaBorderColor: '#d9534f'
          })
          setTimeout(() => {
            this.refresh();
          }, 1000);

        })
      }
      if(this.state.num < this.state.checkNum) {
        this.createPoint(this.getMousePos(e))
        let num = this.state.num;
        ++num;
        this.setState({num: num})
        // this.setState({
        //   num: this.state.num++
        // })
      } 
    }
  }
   //获取坐标
  getMousePos =(e) => {
    const x = e.nativeEvent.offsetX
    const y = e.nativeEvent.offsetY
    return {x, y}
  }
  // 创建坐标点
  createPoint = () => {
    let num = this.state.num
    num ++
    this.setState({
      num:num
    })
  }

  //坐标转换函数
  pointTransfrom = (pointArr,imgSize) => {
    const newPointArr = pointArr.map(p=>{
        let x = Math.round(310 * p.x/parseInt(imgSize.imgWidth)) 
        let y =Math.round(155 * p.y/parseInt(imgSize.imgHeight)) 
        return {x,y}
    })
    return newPointArr
  }
  
  closeBox = () => {
      this.props.verifyPointFixedChild("close")
  }
  
  render() {
    let tempPoints = this.state.tempPoints
    const { vSpace, imgSize, barSize, setSize, isPointShow  } = this.props;
    return (
      // 蒙层
      <div className='mask' style={{ display: isPointShow ? 'block' : 'none' }}>
      <div className='verifybox' style={{ maxWidth: parseInt(imgSize.width) + 30 + 'px' }}>
        <div className='verifybox-top'>
          请完成安全验证
          <span className='verifybox-close' onClick={() => this.closeBox()}>
             <i className='iconfont icon-close'></i>
          </span>
        </div>
        <div className='verifybox-bottom' style={{padding:'15px'}}>
          {/* 验证容器 */}
          <div style={{ position: 'relative' }}>
            <div className='verify-img-out'>
              <div
                className='verify-img-panel'
                style={{
                  width: setSize.imgWidth + 'px',
                  height: setSize.imgHeight + 'px',
                  backgroundSize: setSize.imgWidth + 'px ' + setSize.imgHeight + 'px',
                  marginBottom: vSpace + 'px',
                }}
              >
                <div className='verify-refresh' style={{ zIndex: 3 }} onClick={this.refresh}>
                  <i className='iconfont icon-refresh'></i>
                </div>
                {this.state.pointBackImgBase?
                    <img src={'data:image/png;base64,' + this.state.pointBackImgBase} alt="" style={{width:'100%',height:'100%',display:'block'}} onClick={($event) => this.canvasClick($event)}/>:
                      <img src={defaultImg} alt="" style={{width:'100%',height:'100%',display:'block'}}/>
                }
                {tempPoints.map((tempPoint, index) => {
                  return (
                    <div
                      key={index}
                      className="point-area"
                      style={{
                        backgroundColor: '#1abd6c',
                        color: '#fff',
                        zIndex: 9999,
                        width: '20px',
                        height: '20px',
                        textAlign: 'center',
                        lineHeight: '20px',
                        borderRadius: '50%',
                        position: 'absolute',
                        top: parseInt(tempPoint.y - 10) + 'px',
                        left: parseInt(tempPoint.x - 10) + 'px',
                        overflow:'hidden'
                      }}
                    >{index + 1}</div>
                  );
                })}
              </div>
            </div>

            <div
              className='verify-bar-area'
              style={{
                width: setSize.imgWidth,
                color: this.state.barAreaColor,
                borderColor: this.state.barAreaBorderColor,
                lineHeight: barSize.height,
              }}
            >
              <span className='verify-msg'>{this.state.text}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    );
  }
}

VerifyPointFixed.defaultProps = {
  mode: 'fixed',
  vSpace: 5,
  imgSize: {
    width: '310px',
    height: '200px',
  },
  barSize: {
    width: '310px',
    height: '40px',
  },
  setSize: {
    imgHeight: 200,
    imgWidth: 310,
    barHeight: 0,
    barWidth: 0,
  },
};

export default VerifyPointFixed
