import React, { useState, useEffect, useRef, useLayoutEffect } from 'react';
import styled from "@emotion/styled";
interface EditTextProps {
    text: string;
    canEdit?: boolean;
    blurText?: (text: string) => void;
    changeText?: (text: string) => void;
}
function EditText(props: EditTextProps) {
    // 根据span获取宽度
    const witdthRef = useRef<HTMLDivElement>(null);
    const [showText, setShowText] = useState('');
    const [isFocus, setIsFocus] = useState(false);
    const [inputWith, setInputWith] = useState(100);
    const minTitleWidth = 70;
    const maxTitleWidth = 500;
    useEffect(() => {
        setShowText(props.text);
    }, [props.text]);
    useLayoutEffect(() => {
        dealInputWidth();
    }, [showText]);
    const dealInputWidth = () => {
        const offsetWidth = witdthRef?.current?.offsetWidth || minTitleWidth;
        // +5 防止出现 截断
        const width = offsetWidth < maxTitleWidth ? offsetWidth + 5 : maxTitleWidth;
        setInputWith(width);
    };
    const titleFocus = () => {
        setIsFocus(true);
    };
    const titleInput = (e: React.ChangeEvent<HTMLInputElement>) => {
        const newTitle = e.target.value;
        if (props?.changeText){
            props?.changeText(newTitle);
        }
        setShowText(newTitle);
    };
    const titleBlur = () => {
        const newTitle = showText || '无标题';
        const oldTitle = props.text;
        setIsFocus(false);
        if (showText !== oldTitle) {
            setShowText(newTitle);
            setIsFocus(false);
        } else {
            setIsFocus(false);
            setShowText(newTitle);
        }
        if (props?.blurText) {
            props.blurText(newTitle);
        }
    };
    return (
        <div className='wrap'>
            {props.canEdit ? (
                <EditInput
                    value={showText}
                    style={{ width: inputWith }}
                    onFocus={titleFocus}
                    onChange={titleInput}
                    onBlur={titleBlur}
                    className='input'
                    placeholder="无标题"
                />
            ) : (
                ''
            )}
            {/* 为了计算文字的宽度 */}
            <EtScaleSpan ref={witdthRef}  className={props.canEdit ? 'width' : 'text'} >
              {showText}
            </EtScaleSpan>
        </div>
    );
}

export default EditText

// font-weight: 400;
// white-space: pre;
// visibility: hidden;
// pointer-events: none;
// position: absolute;

const EditInput = styled.input`
  all: unset;
  min-width: 50px;
  margin: 0;
  font-size: var(--titleFontSize);
  font-weight: 400;
  padding: 0 5px 0 0;
  color: #1f1f1f;
  border: none;
`

const EtScaleSpan = styled.span`
  font-weight: 400;
  white-space: pre;
  visibility: hidden;
  pointer-events: none;
  position: absolute;
`
