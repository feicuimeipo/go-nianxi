import React, {useCallback, useState} from "react";
import { ErrorBoundary } from "@/components/error-boundary";
import { FullPageErrorFallback, FullPageLoading } from "@/components/lib";
import {useAuth} from "@/context/auth-context";
import {Route, Routes, useLocation} from "react-router";
import {ServiceAgreement} from "@/unauthenticated-app/agreement/service-agreement";
import {CustomerAgreement} from "@/unauthenticated-app/agreement/customer-agreement";
import {LoginScreen} from "@/unauthenticated-app/login";
import {RegisterScreen} from "@/unauthenticated-app/register";
import {ForgetPasswordScreen} from "@/unauthenticated-app/forget-password";
import {AccountInfoScreen} from "@/page/workstation/view/user-info";
import {useNavigate} from "react-router-dom";
import {useMount} from "@/utils";

const UnAuthenticatedApp = React.lazy(() => import("@/unauthenticated-app"));
const WorkstationApp = React.lazy(() => import("@/page/workstation"));
const ignoredPaths:string[] = ["/serviceAgreement","/serviceAgreement"]



function App() {
    const { user } = useAuth();
    const navigate = useNavigate()
    const location = useLocation()

    const navigateTo = (path:string):string =>{
        ignoredPaths.forEach(v => {
            if (v.startsWith(path)) {
                return path
            }
        })

        if (path==="/auth/login" && user){
            return "/workstation/accountInfo"
        }
        if (path.startsWith("/workstation") && !user){
            return "/auth/login"
        }
        if (path==="" || path==="/"){
            return user?"/workstation/accountInfo":"/auth/login"
        }else{
            return path
        }
    }


      useMount(
          useCallback(() => {
              navigate(navigateTo(location.pathname))
          }, [])
      );

  return (
      <>
        <div className="App">
            <ErrorBoundary fallbackRender={FullPageErrorFallback}>
            <React.Suspense fallback={<FullPageLoading />}>
                <Routes>
                    <Route path={"/workstation"} element={<WorkstationApp   />} >
                        <Route path={"accountInfo"}  element={<AccountInfoScreen />} />
                        <Route path={"useLog"}       element={<AccountInfoScreen  />} />
                        <Route path={"orderInfo"}     element={<AccountInfoScreen />} />
                    </Route>
                    <Route path={"/serviceAgreement"}  element={<ServiceAgreement />} />
                    <Route path={"/customerAgreement"}  element={<CustomerAgreement />} />
                    <Route path={"/auth"}  element={<UnAuthenticatedApp  />} >
                        <Route path={"login"}  element={<LoginScreen />} />
                        <Route path={"register"}  element={<RegisterScreen />} />
                        <Route path={"forgetPassword"}  element={<ForgetPasswordScreen/>} />
                    </Route>
                </Routes>
            </React.Suspense>
          </ErrorBoundary>
        </div>
      </>
  );
}

export default App;
