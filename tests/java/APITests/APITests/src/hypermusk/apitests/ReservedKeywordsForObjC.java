package hypermusk.apitests;





import com.google.gson.annotations.SerializedName;
import java.io.Serializable;

public class ReservedKeywordsForObjC implements Serializable {

	@SerializedName("New")
	private String _new;

	@SerializedName("Alloc")
	private String _alloc;

	@SerializedName("Copy")
	private String _copy;

	@SerializedName("MutableCopy")
	private String _mutableCopy;

	@SerializedName("Description")
	private String _description;

	@SerializedName("NormalName")
	private String _normalName;



	public String getNew() {
		return this._new;
	}
	public void setNew(String _new) {
		this._new = _new;
	}

	public String getAlloc() {
		return this._alloc;
	}
	public void setAlloc(String _alloc) {
		this._alloc = _alloc;
	}

	public String getCopy() {
		return this._copy;
	}
	public void setCopy(String _copy) {
		this._copy = _copy;
	}

	public String getMutableCopy() {
		return this._mutableCopy;
	}
	public void setMutableCopy(String _mutableCopy) {
		this._mutableCopy = _mutableCopy;
	}

	public String getDescription() {
		return this._description;
	}
	public void setDescription(String _description) {
		this._description = _description;
	}

	public String getNormalName() {
		return this._normalName;
	}
	public void setNormalName(String _normalName) {
		this._normalName = _normalName;
	}



}

