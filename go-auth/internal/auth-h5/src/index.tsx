import "./wdyr";
import React from "react";
import ReactDOM from 'react-dom/client';
import App from "./App";
import {ConfigProvider} from "antd";
import { AppProviders } from "@/context";
import { Profiler } from "@/components/profiler";


const root = ReactDOM.createRoot(
    document.getElementById('root') as HTMLElement
);


root.render(
    <React.StrictMode>
        <Profiler id={"root"} phases={["mount"]}>
            <ConfigProvider  theme={{
                token: {
                    colorPrimary: '#00b96b',
                    colorLink: '#00b96b',
                    colorLinkHover: '#00B96BB6',
                }
            }}
            >
                <AppProviders>
                  <App />
                </AppProviders>
            </ConfigProvider>
        </Profiler>
    </React.StrictMode>
)



