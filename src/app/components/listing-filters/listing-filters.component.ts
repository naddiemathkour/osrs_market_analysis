import {
  Component,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
} from '@angular/core';
import {
  MatButtonToggleChange,
  MatButtonToggleModule,
} from '@angular/material/button-toggle';
import { MatSliderModule } from '@angular/material/slider';
import { MatSelectModule } from '@angular/material/select';
import {
  MatCheckboxChange,
  MatCheckboxModule,
} from '@angular/material/checkbox';
import { IFilters } from '../../interfaces/filters.interface';

@Component({
  selector: 'app-listing-filters',
  standalone: true,
  imports: [
    MatButtonToggleModule,
    MatSliderModule,
    MatSelectModule,
    MatCheckboxModule,
  ],
  templateUrl: './listing-filters.component.html',
  styleUrl: './listing-filters.component.scss',
})
export class ListingFiltersComponent implements OnChanges, OnInit {
  @Input() filterValues!: IFilters;
  selectedFilter: string | null = null;
  filterNames = {
    buyLimit: 'Buy Limit',
    margin: 'Profit Margin',
    alchprof: 'Alch Profit',
    highVolume: 'High Volume',
    lowVolume: 'Low Volume',
  };

  ngOnInit(): void {}

  ngOnChanges(changes: SimpleChanges): void {
    console.log(changes);
  }

  toggleView(event: MatButtonToggleChange) {
    this.filterValues.dataType = event.value;
  }

  toggleMembers(event: MatCheckboxChange) {
    this.filterValues.members = event.checked;
    console.log(this.filterValues);
  }
}
