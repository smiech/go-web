import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { LgbuttonComponent } from './lgbutton.component';

describe('LgbuttonComponent', () => {
  let component: LgbuttonComponent;
  let fixture: ComponentFixture<LgbuttonComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ LgbuttonComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LgbuttonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
