import { Component, OnInit } from '@angular/core';
import { PersistenceService } from '../../services/persistence.service';
import { ActivatedRoute } from '@angular/router';
import { Http } from '@angular/http';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  constructor(private route: ActivatedRoute,
    private presistenceService: PersistenceService,
    private authService: AuthService,
    private http: Http,
    private router: Router) {
  }

  ngOnInit() {
    if (!this.authService.checkLogin()) {
      this.router.navigate([environment.urls.login]);
    }
    var authHeaders = this.authService.initAuthHeaders();
  }
}
