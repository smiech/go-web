import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserspaceComponent } from './userspace.component';

describe('UserspaceComponent', () => {
  let component: UserspaceComponent;
  let fixture: ComponentFixture<UserspaceComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserspaceComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserspaceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
