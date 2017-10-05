import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './components/app/app.component';
import { LoginComponent } from './components/login/login.component';
import { NavMenuComponent } from './components/navmenu/navmenu.component';

import { AuthService } from './services/auth.service';
import { PersistenceService } from './services/persistence.service';
import { HomeComponent } from './components/home/home.component';
import { UserspaceComponent } from './components/userspace/userspace.component';
import { environment } from '../environments/environment';
import { LgbuttonComponent } from './components/lgbutton/lgbutton.component';
import { Lgbutton0Component } from './components/lgbutton0/lgbutton0.component';
import { StyleguideComponent } from './components/styleguide/styleguide.component';
import { LgButtonComponent } from './components/lg-button/lg-button.component';
import { LgFormFieldComponent } from './components/lg-form-field/lg-form-field.component'

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    NavMenuComponent,
    HomeComponent,
    UserspaceComponent,
    LgbuttonComponent,
    Lgbutton0Component,
    StyleguideComponent,
    LgButtonComponent,
    LgFormFieldComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    RouterModule.forRoot([
      { path: '', redirectTo: environment.urls.home, pathMatch: 'full' },
      { path: environment.urls.home, component: HomeComponent },
      { path: environment.urls.login, component: LoginComponent },
      { path: environment.urls.userspace, component: UserspaceComponent },
      { path: environment.urls.styleguide, component: StyleguideComponent },
    ]),

  ],
  providers: [AuthService, PersistenceService],
  bootstrap: [AppComponent]
})
export class AppModule { }
