import { NgModule } from '@angular/core';

import { SharedModule } from '@shared';
import { UserActiveComponent } from './account/active/active.component';
// account pages
import { UserLoginComponent } from './account/login/login.component';
import { UserRegisterResultComponent } from './account/register-result/register-result.component';
import { UserRegisterComponent } from './account/register/register.component';
// single pages
import { CallbackComponent } from './callback/callback.component';
// dashboard pages
import { DashboardComponent } from './dashboard/dashboard.component';
import { RouteRoutingModule } from './routes-routing.module';

const COMPONENTS = [
  DashboardComponent,
  // account pages
  UserLoginComponent,
  UserRegisterComponent,
  UserRegisterResultComponent,
  UserActiveComponent,
  // single pages
  CallbackComponent,
];
const COMPONENTS_NOROUNT = [];

@NgModule({
  imports: [ SharedModule, RouteRoutingModule ],
  declarations: [
    ...COMPONENTS,
    ...COMPONENTS_NOROUNT
  ],
})
export class RoutesModule {}
