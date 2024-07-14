import {
  Component,
  Input,
  OnChanges,
  OnInit,
  EventEmitter,
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
export class ListingFiltersComponent implements OnInit {
  @Input() filterValues!: IFilters;
  @Output() updateFilter: EventEmitter<IFilters> = new EventEmitter<IFilters>();
  newFilters!: IFilters;
  selectedFilter: string | null = null;

  filterNames = {
    buyLimit: 'Buy Limit',
    margin: 'Profit Margin',
    alchprof: 'Alch Profit',
    highVolume: 'High Volume',
    lowVolume: 'Low Volume',
  };

  ngOnInit(): void {
    this.newFilters = {
      dataType: this.filterValues.dataType,
      alchprof: { max: 0, filter: 1 },
      margin: { max: 0, filter: 1 },
      buyLimit: { max: 0, filter: 1 },
      highVolume: { max: 0, filter: 1 },
      lowVolume: { max: 0, filter: 1 },
      members: true,
    };
  }

  toggleView(event: MatButtonToggleChange) {
    this.newFilters.dataType = event.value;
  }

  toggleMembers(event: MatCheckboxChange) {
    this.newFilters.members = !event.checked;
  }

  setFilter(event: any) {
    const filterVal = Number(event.target.value);
    switch (this.selectedFilter) {
      case this.filterNames.alchprof:
        this.newFilters.alchprof.filter = filterVal;
        break;
      case this.filterNames.margin:
        this.newFilters.margin.filter = filterVal;
        break;
      case this.filterNames.buyLimit:
        this.newFilters.buyLimit.filter = filterVal;
        break;
      case this.filterNames.highVolume:
        this.newFilters.highVolume.filter = filterVal;
        break;
      case this.filterNames.lowVolume:
        this.newFilters.lowVolume.filter = filterVal;
        break;

      default:
        console.error(
          'Error saving filter. Filter field name does not exist...'
        );
        break;
    }
  }

  onFilter(): void {
    this.updateFilter.emit(this.newFilters);
  }
}
