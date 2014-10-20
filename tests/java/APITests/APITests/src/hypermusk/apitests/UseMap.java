package hypermusk.apitests;




import java.util.Map;
import com.google.gson.annotations.SerializedName;
import java.io.Serializable;

public class UseMap implements Serializable {

	@SerializedName("Map")
	private Map _map;



	public Map getMap() {
		return this._map;
	}
	public void setMap(Map _map) {
		this._map = _map;
	}



}

