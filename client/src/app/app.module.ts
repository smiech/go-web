import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { AppComponent } from './components/app/app.component';
import { LoginComponent } from './components/login/login.component';

import { AuthService } from './services/auth.service';
import { PersistenceService } from './services/persistence.service';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent
  ],
  imports: [
    BrowserModule,
    FormsModule,
    HttpModule,
    RouterModule.forRoot([
      { path: 'login', component: LoginComponent },
      { path: 'login/:email', component: LoginComponent },
    ]),
   
  ],
  providers: [AuthService, PersistenceService],
  bootstrap: [AppComponent]
})
export class AppModule { }
