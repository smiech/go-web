import { Component, OnInit } from '@angular/core';
import { PersistenceService } from '../../services/persistence.service';
import { ActivatedRoute } from '@angular/router';
import { Http } from '@angular/http';
import { Router } from '@angular/router';
import { AuthService } from '../../services/auth.service';
import { DeviceService } from '../../services/device.service';
import { environment } from '../../../environments/environment';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {

  constructor(private route: ActivatedRoute,
    private persistenceService: PersistenceService,
    private authService: AuthService,
    private deviceService: DeviceService,
    private http: Http,
    private router: Router) {
  }

  ngOnInit() {
    if (!this.authService.checkLogin()) {
      this.router.navigate([environment.urls.login]);
    }

    this.getDevices();
   
    var authHeaders = this.authService.initAuthHeaders();
  }
  private getDevices(): void {
    this.deviceService
        .getDevices(this.persistenceService.currentUser.id)
        .then(x => {
            // this.currentMeasure = x;
            // this.currentMeasure.description = this.currentMeasure.description.replace(new RegExp('\r?\n', 'g'), '<br />');
        })
        .catch(err => {
            console.log(err);
        }); 
}
}
