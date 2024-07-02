import React,{CSSProperties} from 'react'
import * as icons  from '@ant-design/icons'

const Icon = (props:{ icon: string,style?: CSSProperties}) => {
    const { icon,style } = props;
    const antIcon: { [key: string]: any } = icons;
    if (props.style){
        antIcon[icon].style = props.style
    }
    return style? React.createElement(antIcon[icon],{style: style}):
        React.createElement(antIcon[icon]);
};


export default Icon
